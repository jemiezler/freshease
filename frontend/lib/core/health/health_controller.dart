// lib/core/health/health_controller.dart
import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:health/health.dart';
import 'package:frontend/core/health/health_repository.dart';
import 'package:frontend/core/health/util.dart'; // <- where you keep dataTypesAndroid/dataTypesIOS

enum HealthState { idle, fetching, ready, noData, authDenied, error }

class HealthController extends ChangeNotifier {
  HealthController({HealthRepository? repository})
    : _repo = repository ?? NoopHealthRepository();

  final Health _health = Health();
  final HealthRepository _repo;

  HealthState state = HealthState.idle;
  HealthConnectSdkStatus? hcStatus;

  // ====== What the page wants ======
  int stepsToday = 0;
  double kcalToday = 0;

  // ====== What you want to save ======
  List<HealthDataPoint> allPoints = [];

  // Types: collect EVERYTHING supported for the platform
  List<HealthDataType> get _allTypes => Platform.isAndroid
      ? dataTypesAndroid // from your util.dart (same as in your sample)
      : (Platform.isIOS ? dataTypesIOS : const []);

  // READ/WRITE where allowed; READ-only for iOS-restricted types is handled below
  List<HealthDataAccess> get _permissions => _allTypes.map((type) {
    const iosReadOnly = <HealthDataType>{
      HealthDataType.GENDER,
      HealthDataType.BLOOD_TYPE,
      HealthDataType.BIRTH_DATE,
      HealthDataType.APPLE_MOVE_TIME,
      HealthDataType.APPLE_STAND_HOUR,
      HealthDataType.APPLE_STAND_TIME,
      HealthDataType.WALKING_HEART_RATE,
      HealthDataType.ELECTROCARDIOGRAM,
      HealthDataType.HIGH_HEART_RATE_EVENT,
      HealthDataType.LOW_HEART_RATE_EVENT,
      HealthDataType.IRREGULAR_HEART_RATE_EVENT,
      HealthDataType.EXERCISE_TIME,
    };
    if (Platform.isIOS && iosReadOnly.contains(type)) {
      return HealthDataAccess.READ;
    }
    return HealthDataAccess.READ_WRITE;
  }).toList();

  Future<void> init() async {
    _health.configure();
    if (Platform.isAndroid) {
      hcStatus = await _health.getHealthConnectSdkStatus();
    }
    notifyListeners();
  }

  Future<void> authorize() async {
    try {
      if (Platform.isAndroid) {
        await Permission.activityRecognition.request();
        await Permission.location.request(); // some calories/workouts gate this
      }

      bool? has = await _health.hasPermissions(
        _allTypes,
        permissions: _permissions,
      );
      has = has == true;

      if (!has) {
        final ok = await _health.requestAuthorization(
          _allTypes,
          permissions: _permissions,
        );
        if (!ok) {
          state = HealthState.authDenied;
          notifyListeners();
          return;
        }
        // optional (Android HC)
        await _health.requestHealthDataHistoryAuthorization().catchError(
          (_) {},
        );
        await _health.requestHealthDataInBackgroundAuthorization().catchError(
          (_) {},
        );
      }
    } catch (_) {
      state = HealthState.authDenied;
      notifyListeners();
    }
  }

  /// Fetches ALL points for last 24h (saved) + computes only two KPIs for UI.
  Future<void> fetchAll24hAndComputeKpis() async {
    state = HealthState.fetching;
    notifyListeners();

    final now = DateTime.now();
    final start = now.subtract(const Duration(hours: 24));

    try {
      // 1) Fetch ALL points
      final fetched = await _health.getHealthDataFromTypes(
        types: _allTypes,
        startTime: start,
        endTime: now,
      );

      // dedup
      final dedup = _health.removeDuplicates(fetched);

      // keep in memory
      allPoints = List.unmodifiable(dedup);

      // 2) Persist (optional)
      await _repo.savePoints(allPoints);

      // 3) Compute KPIs for the page
      // Steps (use aggregate helper to match your logs/output)
      final int? s = await _health.getTotalStepsInInterval(
        DateTime(now.year, now.month, now.day),
        now,
        includeManualEntry: true,
      );
      stepsToday = s ?? 0;

      // Calories: sum TOTAL_CALORIES_BURNED points (kcal)
      double sumKcal = 0;
      for (final p in allPoints) {
        if (p.type == HealthDataType.TOTAL_CALORIES_BURNED) {
          final v = p.value;
          if (v is num) {
            final parsed = double.tryParse(v.toString());
            if (parsed != null) sumKcal += parsed;
          } else {
            final parsed = double.tryParse(v.toString());
            if (parsed != null) sumKcal += parsed;
          }
        }
      }
      kcalToday = sumKcal;

      state = allPoints.isEmpty ? HealthState.noData : HealthState.ready;
    } catch (e) {
      state = HealthState.error;
    } finally {
      notifyListeners();
    }
  }
}

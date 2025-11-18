import 'package:flutter/foundation.dart';
import 'package:frontend/core/health/util.dart';
import 'package:frontend/core/platform_helper.dart';

// Conditional imports for mobile-only packages
import 'package:health/health.dart'
    if (dart.library.html) 'package:frontend/core/health/health_stub.dart'
    as health_pkg;
import 'package:permission_handler/permission_handler.dart'
    if (dart.library.html) 'package:frontend/core/health/permission_stub.dart'
    as permission_pkg;

class HealthService {
  HealthService._();
  static final HealthService instance = HealthService._();

  final health_pkg.Health _health = health_pkg.Health();
  health_pkg.HealthConnectSdkStatus? _hcStatus;

  // Cached
  List<health_pkg.HealthDataPoint> _lastData = [];
  int _lastSteps = 0;

  List<health_pkg.HealthDataPoint> get lastData => _lastData;
  int get lastSteps => _lastSteps;
  health_pkg.HealthConnectSdkStatus? get healthConnectStatus => _hcStatus;

  /// Data types you want to work with (platform-aware)
  List<health_pkg.HealthDataType> get types {
    if (kIsWeb) return const [];
    if (PlatformHelper.isAndroid) {
      return dataTypesAndroid;
    } else if (PlatformHelper.isIOS) {
      return dataTypesIOS;
    }
    return const [];
  }

  /// READ / WRITE matrix (iOS has READ-only for some types)
  List<health_pkg.HealthDataAccess> get permissions => types
      .map(
        (type) =>
            const {
              health_pkg.HealthDataType.GENDER,
              health_pkg.HealthDataType.BLOOD_TYPE,
              health_pkg.HealthDataType.BIRTH_DATE,
              health_pkg.HealthDataType.APPLE_MOVE_TIME,
              health_pkg.HealthDataType.APPLE_STAND_HOUR,
              health_pkg.HealthDataType.APPLE_STAND_TIME,
              health_pkg.HealthDataType.WALKING_HEART_RATE,
              health_pkg.HealthDataType.ELECTROCARDIOGRAM,
              health_pkg.HealthDataType.HIGH_HEART_RATE_EVENT,
              health_pkg.HealthDataType.LOW_HEART_RATE_EVENT,
              health_pkg.HealthDataType.IRREGULAR_HEART_RATE_EVENT,
              health_pkg.HealthDataType.EXERCISE_TIME,
            }.contains(type)
            ? health_pkg.HealthDataAccess.READ
            : health_pkg.HealthDataAccess.READ_WRITE,
      )
      .toList();

  /// Call this once at startup (before you use any methods).
  Future<void> init() async {
    if (kIsWeb) return;
    _health.configure();

    if (PlatformHelper.isAndroid) {
      _hcStatus = await _health.getHealthConnectSdkStatus();
    }
  }

  /// (Android) deep link to HC install
  Future<void> installHealthConnect() async {
    if (kIsWeb) return;
    if (!PlatformHelper.isAndroid) return;
    await _health.installHealthConnect();
  }

  Future<bool> authorize() async {
    if (kIsWeb) return false;

    // For steps/workouts/sleep etc.
    await permission_pkg.Permission.activityRecognition.request();
    // For distance in workouts
    await permission_pkg.Permission.location.request();

    try {
      final ok = await _health.requestAuthorization(
        types,
        permissions: permissions,
      );
      await _health.requestHealthDataHistoryAuthorization();
      await _health.requestHealthDataInBackgroundAuthorization();
      return ok;
    } catch (e) {
      debugPrint('[HealthService] authorize error: $e');
      return false;
    }
  }

  Future<List<health_pkg.HealthDataPoint>> fetchLatest24h({
    List<health_pkg.RecordingMethod> exclude = const [],
  }) async {
    if (kIsWeb) {
      _lastData = [];
      return _lastData;
    }

    final now = DateTime.now();
    final yesterday = now.subtract(const Duration(hours: 24));
    try {
      final points = await _health.getHealthDataFromTypes(
        types: types,
        startTime: yesterday,
        endTime: now,
        recordingMethodsToFilter: exclude,
      );
      points.sort((a, b) => b.dateTo.compareTo(a.dateTo));
      _lastData = _health.removeDuplicates(
        points.length <= 100 ? points : points.sublist(0, 100),
      );
    } catch (e) {
      debugPrint('[HealthService] fetchLatest24h error: $e');
      _lastData = [];
    }
    return _lastData;
  }

  Future<health_pkg.HealthDataPoint?> fetchByUUID({
    required String uuid,
    required health_pkg.HealthDataType type,
  }) async {
    if (kIsWeb) return null;
    try {
      return await _health.getHealthDataByUUID(uuid: uuid, type: type);
    } catch (e) {
      debugPrint('[HealthService] fetchByUUID error: $e');
      return null;
    }
  }

  Future<int> fetchTodaySteps({bool includeManual = true}) async {
    if (kIsWeb) {
      _lastSteps = 0;
      return _lastSteps;
    }

    final now = DateTime.now();
    final midnight = DateTime(now.year, now.month, now.day);

    bool ok =
        await _health.hasPermissions([health_pkg.HealthDataType.STEPS]) ??
        false;
    if (!ok)
      ok = await _health.requestAuthorization([
        health_pkg.HealthDataType.STEPS,
      ]);

    if (!ok) {
      _lastSteps = 0;
      return _lastSteps;
    }

    try {
      final steps = await _health.getTotalStepsInInterval(
        midnight,
        now,
        includeManualEntry: includeManual,
      );
      _lastSteps = steps ?? 0;
    } catch (e) {
      debugPrint('[HealthService] fetchTodaySteps error: $e');
      _lastSteps = 0;
    }
    return _lastSteps;
  }

  Future<bool> revokePermissions() async {
    if (kIsWeb) return false;
    try {
      await _health.revokePermissions();
      return true;
    } catch (e) {
      debugPrint('[HealthService] revoke error: $e');
      return false;
    }
  }

  /// Examples of writing/deleting (trim to your needs)
  Future<bool> writeSample({
    required double value,
    required health_pkg.HealthDataType type,
    required DateTime start,
    DateTime? end,
    health_pkg.RecordingMethod recordingMethod =
        health_pkg.RecordingMethod.manual,
  }) async {
    if (kIsWeb) return false;
    try {
      return await _health.writeHealthData(
        value: value,
        type: type,
        startTime: start,
        endTime: end ?? start,
        recordingMethod: recordingMethod,
      );
    } catch (e) {
      debugPrint('[HealthService] writeSample error: $e');
      return false;
    }
  }

  Future<bool> deleteRange({
    required health_pkg.HealthDataType type,
    required DateTime start,
    required DateTime end,
  }) async {
    if (kIsWeb) return false;
    try {
      return await _health.delete(type: type, startTime: start, endTime: end);
    } catch (e) {
      debugPrint('[HealthService] deleteRange error: $e');
      return false;
    }
  }
}

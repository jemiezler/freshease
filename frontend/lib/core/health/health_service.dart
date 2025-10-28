import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:frontend/core/health/util.dart';
import 'package:health/health.dart';
import 'package:permission_handler/permission_handler.dart';

class HealthService {
  HealthService._();
  static final HealthService instance = HealthService._();

  final Health _health = Health();
  HealthConnectSdkStatus? _hcStatus;

  // Cached
  List<HealthDataPoint> _lastData = [];
  int _lastSteps = 0;

  List<HealthDataPoint> get lastData => _lastData;
  int get lastSteps => _lastSteps;
  HealthConnectSdkStatus? get healthConnectStatus => _hcStatus;

  /// Data types you want to work with (platform-aware)
  List<HealthDataType> get types => Platform.isAndroid
      ? dataTypesAndroid
      : Platform.isIOS
      ? dataTypesIOS
      : const [];

  /// READ / WRITE matrix (iOS has READ-only for some types)
  List<HealthDataAccess> get permissions => types
      .map(
        (type) =>
            const {
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
            }.contains(type)
            ? HealthDataAccess.READ
            : HealthDataAccess.READ_WRITE,
      )
      .toList();

  /// Call this once at startup (before you use any methods).
  Future<void> init() async {
    _health.configure();

    if (Platform.isAndroid) {
      _hcStatus = await _health.getHealthConnectSdkStatus();
    }
  }

  /// (Android) deep link to HC install
  Future<void> installHealthConnect() async {
    if (!Platform.isAndroid) return;
    await _health.installHealthConnect();
  }

  Future<bool> authorize() async {
    // For steps/workouts/sleep etc.
    await Permission.activityRecognition.request();
    // For distance in workouts
    await Permission.location.request();

    // hasPermissions returns false if WRITE can’t be disclosed—force re-request
    bool? has = await _health.hasPermissions(types, permissions: permissions);
    has = false;

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

  Future<List<HealthDataPoint>> fetchLatest24h({
    List<RecordingMethod> exclude = const [],
  }) async {
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

  Future<HealthDataPoint?> fetchByUUID({
    required String uuid,
    required HealthDataType type,
  }) async {
    try {
      return await _health.getHealthDataByUUID(uuid: uuid, type: type);
    } catch (e) {
      debugPrint('[HealthService] fetchByUUID error: $e');
      return null;
    }
  }

  Future<int> fetchTodaySteps({bool includeManual = true}) async {
    final now = DateTime.now();
    final midnight = DateTime(now.year, now.month, now.day);

    bool ok = await _health.hasPermissions([HealthDataType.STEPS]) ?? false;
    if (!ok) ok = await _health.requestAuthorization([HealthDataType.STEPS]);

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
    required HealthDataType type,
    required DateTime start,
    DateTime? end,
    RecordingMethod recordingMethod = RecordingMethod.manual,
  }) async {
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
    required HealthDataType type,
    required DateTime start,
    required DateTime end,
  }) async {
    try {
      return await _health.delete(type: type, startTime: start, endTime: end);
    } catch (e) {
      debugPrint('[HealthService] deleteRange error: $e');
      return false;
    }
  }
}

// Stub implementation for web platform
// Note: This file is only used when compiling for web (dart.library.html exists)
// The actual health package types are imported conditionally in util.dart

// These enums must match the health package enums
enum HealthConnectSdkStatus { installed, notInstalled, notSupported }

enum HealthDataType {
  ACTIVE_ENERGY_BURNED,
  APPLE_STAND_TIME,
  APPLE_STAND_HOUR,
  APPLE_MOVE_TIME,
  AUDIOGRAM,
  BASAL_ENERGY_BURNED,
  BLOOD_GLUCOSE,
  BLOOD_OXYGEN,
  BLOOD_PRESSURE_DIASTOLIC,
  BLOOD_PRESSURE_SYSTOLIC,
  BODY_FAT_PERCENTAGE,
  BODY_MASS_INDEX,
  BODY_TEMPERATURE,
  DIETARY_CARBS_CONSUMED,
  DIETARY_CAFFEINE,
  DIETARY_ENERGY_CONSUMED,
  DIETARY_FATS_CONSUMED,
  DIETARY_PROTEIN_CONSUMED,
  ELECTRODERMAL_ACTIVITY,
  FORCED_EXPIRATORY_VOLUME,
  HEART_RATE,
  HEART_RATE_VARIABILITY_SDNN,
  HEIGHT,
  INSULIN_DELIVERY,
  RESPIRATORY_RATE,
  PERIPHERAL_PERFUSION_INDEX,
  STEPS,
  WAIST_CIRCUMFERENCE,
  WEIGHT,
  FLIGHTS_CLIMBED,
  DISTANCE_WALKING_RUNNING,
  WALKING_SPEED,
  MINDFULNESS,
  SLEEP_AWAKE,
  SLEEP_ASLEEP,
  SLEEP_IN_BED,
  SLEEP_LIGHT,
  SLEEP_DEEP,
  SLEEP_REM,
  WATER,
  EXERCISE_TIME,
  WORKOUT,
  HEADACHE_NOT_PRESENT,
  HEADACHE_MILD,
  HEADACHE_MODERATE,
  HEADACHE_SEVERE,
  HEADACHE_UNSPECIFIED,
  LEAN_BODY_MASS,
  ELECTROCARDIOGRAM,
  NUTRITION,
  GENDER,
  BLOOD_TYPE,
  BIRTH_DATE,
  MENSTRUATION_FLOW,
  WATER_TEMPERATURE,
  UNDERWATER_DEPTH,
  UV_INDEX,
  HEART_RATE_VARIABILITY_RMSSD,
  DISTANCE_DELTA,
  SPEED,
  SLEEP_AWAKE_IN_BED,
  SLEEP_OUT_OF_BED,
  SLEEP_UNKNOWN,
  SLEEP_SESSION,
  RESTING_HEART_RATE,
  TOTAL_CALORIES_BURNED,
  WALKING_HEART_RATE,
  HIGH_HEART_RATE_EVENT,
  LOW_HEART_RATE_EVENT,
  IRREGULAR_HEART_RATE_EVENT,
}

enum HealthDataAccess { READ, READ_WRITE }

enum RecordingMethod { manual, automatic }

class HealthDataPoint {
  final HealthDataType type;
  final dynamic value;
  final DateTime dateTo;
  HealthDataPoint({
    required this.type,
    required this.value,
    required this.dateTo,
  });
}

class Health {
  void configure() {}
  Future<HealthConnectSdkStatus?> getHealthConnectSdkStatus() async => null;
  Future<bool> requestAuthorization(
    List<HealthDataType> types, {
    List<HealthDataAccess>? permissions,
  }) async => false;
  Future<bool?> hasPermissions(
    List<HealthDataType> types, {
    List<HealthDataAccess>? permissions,
  }) async => false;
  Future<void> requestHealthDataHistoryAuthorization() async {}
  Future<void> requestHealthDataInBackgroundAuthorization() async {}
  Future<List<HealthDataPoint>> getHealthDataFromTypes({
    required List<HealthDataType> types,
    required DateTime startTime,
    required DateTime endTime,
    List<RecordingMethod>? recordingMethodsToFilter,
  }) async => [];
  List<HealthDataPoint> removeDuplicates(List<HealthDataPoint> points) =>
      points;
  Future<HealthDataPoint?> getHealthDataByUUID({
    required String uuid,
    required HealthDataType type,
  }) async => null;
  Future<int?> getTotalStepsInInterval(
    DateTime start,
    DateTime end, {
    bool includeManualEntry = true,
  }) async => 0;
  Future<bool> writeHealthData({
    required double value,
    required HealthDataType type,
    required DateTime startTime,
    required DateTime endTime,
    RecordingMethod recordingMethod = RecordingMethod.manual,
  }) async => false;
  Future<bool> delete({
    required HealthDataType type,
    required DateTime startTime,
    required DateTime endTime,
  }) async => false;
  Future<void> revokePermissions() async {}
  Future<void> installHealthConnect() async {}
}

import 'package:health/health.dart';

abstract class HealthRepository {
  Future<void> savePoints(List<HealthDataPoint> points);
}

// No-op default so you don’t need a DB right away
class NoopHealthRepository implements HealthRepository {
  @override
  Future<void> savePoints(List<HealthDataPoint> points) async {
    /* noop */
  }
}

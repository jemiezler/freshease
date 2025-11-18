// Conditional import for health package
import 'package:health/health.dart'
    if (dart.library.html) 'package:frontend/core/health/health_stub.dart'
    as health_pkg;

abstract class HealthRepository {
  Future<void> savePoints(List<health_pkg.HealthDataPoint> points);
}

// No-op default so you don't need a DB right away
class NoopHealthRepository implements HealthRepository {
  @override
  Future<void> savePoints(List<health_pkg.HealthDataPoint> points) async {
    /* noop */
  }
}

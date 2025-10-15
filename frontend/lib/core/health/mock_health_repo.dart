import 'dart:async';
import 'dart:math';

/// Data model you can reuse later with a real repo
class CalorieSnapshot {
  final double intakeKcal; // food in (kcal)
  final double activeBurnKcal; // activity out (kcal)
  final DateTime syncedAt;

  CalorieSnapshot({
    required this.intakeKcal,
    required this.activeBurnKcal,
    required this.syncedAt,
  });

  double get netKcal => intakeKcal - activeBurnKcal;
}

/// Interface for future real implementation (Health Connect / HealthKit)
abstract class IHealthRepo {
  Future<CalorieSnapshot> readToday();
}

/// Mock implementation
class MockHealthRepo implements IHealthRepo {
  final _rng = Random();
  @override
  Future<CalorieSnapshot> readToday() async {
    // Simulate latency
    await Future.delayed(const Duration(milliseconds: 600));

    // Generate plausible numbers
    // Intake: 1,600–2,400 kcal, Burn: 400–900 kcal (active)
    final intake = 1600 + _rng.nextInt(800); // 1600..2399
    final burn = 400 + _rng.nextInt(500); // 400..899

    return CalorieSnapshot(
      intakeKcal: intake.toDouble(),
      activeBurnKcal: burn.toDouble(),
      syncedAt: DateTime.now(),
    );
  }
}

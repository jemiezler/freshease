class MealPlan {
  final String day;
  final Map<String, String> meals;
  final Map<String, int> calories;
  final int totalCalories;

  MealPlan({
    required this.day,
    required this.meals,
    required this.calories,
    required this.totalCalories,
  });

  factory MealPlan.fromJson(Map<String, dynamic> json) {
    return MealPlan(
      day: json['day'] as String? ?? '',
      meals: Map<String, String>.from(json['meals'] as Map? ?? {}),
      calories: Map<String, int>.from(json['calories'] as Map? ?? {}),
      totalCalories: json['total_calories'] as int? ?? 0,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'day': day,
      'meals': meals,
      'calories': calories,
      'total_calories': totalCalories,
    };
  }
}

class GenAiResponse {
  final int stepsToday;
  final double activeKcal24h;
  final List<MealPlan> plan;

  GenAiResponse({
    required this.stepsToday,
    required this.activeKcal24h,
    required this.plan,
  });

  factory GenAiResponse.fromJson(Map<String, dynamic> json) {
    return GenAiResponse(
      stepsToday: json['steps_today'] as int? ?? 0,
      activeKcal24h: (json['active_kcal_24h'] as num?)?.toDouble() ?? 0.0,
      plan:
          (json['plan'] as List?)
              ?.map((item) => MealPlan.fromJson(item as Map<String, dynamic>))
              .toList() ??
          [],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'steps_today': stepsToday,
      'active_kcal_24h': activeKcal24h,
      'plan': plan.map((mealPlan) => mealPlan.toJson()).toList(),
    };
  }
}

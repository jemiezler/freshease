import 'package:flutter/material.dart';
import '../genai/models.dart';

class MealPlanCard extends StatelessWidget {
  final MealPlan mealPlan;

  const MealPlanCard({super.key, required this.mealPlan});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      margin: const EdgeInsets.symmetric(vertical: 8),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  mealPlan.day,
                  style: theme.textTheme.titleLarge?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 12,
                    vertical: 6,
                  ),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.primaryContainer,
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Text(
                    '${mealPlan.totalCalories} kcal',
                    style: theme.textTheme.labelMedium?.copyWith(
                      color: theme.colorScheme.onPrimaryContainer,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            ...mealPlan.meals.entries.map((entry) {
              final mealType = entry.key;
              final mealDescription = entry.value;
              final calories = mealPlan.calories[mealType] ?? 0;

              return _MealItem(
                mealType: mealType,
                description: mealDescription,
                calories: calories,
              );
            }),
          ],
        ),
      ),
    );
  }
}

class _MealItem extends StatelessWidget {
  final String mealType;
  final String description;
  final int calories;

  const _MealItem({
    required this.mealType,
    required this.description,
    required this.calories,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 80,
            child: Text(
              mealType.toUpperCase(),
              style: theme.textTheme.labelMedium?.copyWith(
                fontWeight: FontWeight.bold,
                color: theme.colorScheme.primary,
              ),
            ),
          ),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(description, style: theme.textTheme.bodyMedium),
                const SizedBox(height: 4),
                Text(
                  '$calories kcal',
                  style: theme.textTheme.labelSmall?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class MealPlanGenerator extends StatefulWidget {
  final VoidCallback onGenerate;
  final bool isLoading;
  final String? error;
  final String? userGender;
  final int? userAge;
  final String? userGoal;
  final double? userHeightCm;
  final double? userWeightKg;
  final bool showUserData;

  const MealPlanGenerator({
    super.key,
    required this.onGenerate,
    required this.isLoading,
    this.error,
    this.userGender,
    this.userAge,
    this.userGoal,
    this.userHeightCm,
    this.userWeightKg,
    this.showUserData = true,
  });

  @override
  State<MealPlanGenerator> createState() => _MealPlanGeneratorState();
}

class _MealPlanGeneratorState extends State<MealPlanGenerator> {
  @override
  void initState() {
    super.initState();
    // User data is now passed via widget properties
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      margin: const EdgeInsets.symmetric(vertical: 8),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'AI Meal Plan Generator',
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),

            // Show comprehensive user data if available
            if (widget.showUserData &&
                (widget.userGender != null ||
                    widget.userAge != null ||
                    widget.userGoal != null ||
                    widget.userHeightCm != null ||
                    widget.userWeightKg != null))
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: theme.colorScheme.primaryContainer.withValues(
                    alpha: 0.1,
                  ),
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(
                    color: theme.colorScheme.primary.withValues(alpha: 0.3),
                  ),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Icon(
                          Icons.person,
                          color: theme.colorScheme.primary,
                          size: 20,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          'Using your profile data:',
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color: theme.colorScheme.primary,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Wrap(
                      spacing: 16,
                      runSpacing: 8,
                      children: [
                        if (widget.userGender != null)
                          _ProfileChip(
                            label: 'Gender',
                            value: widget.userGender!,
                            icon: Icons.person_outline,
                          ),
                        if (widget.userAge != null)
                          _ProfileChip(
                            label: 'Age',
                            value: '${widget.userAge} years',
                            icon: Icons.cake,
                          ),
                        if (widget.userGoal != null)
                          _ProfileChip(
                            label: 'Goal',
                            value: widget.userGoal!
                                .replaceAll('_', ' ')
                                .toUpperCase(),
                            icon: Icons.flag,
                          ),
                        if (widget.userHeightCm != null)
                          _ProfileChip(
                            label: 'Height',
                            value:
                                '${widget.userHeightCm!.toStringAsFixed(0)} cm',
                            icon: Icons.height,
                          ),
                        if (widget.userWeightKg != null)
                          _ProfileChip(
                            label: 'Weight',
                            value:
                                '${widget.userWeightKg!.toStringAsFixed(1)} kg',
                            icon: Icons.monitor_weight,
                          ),
                      ],
                    ),
                  ],
                ),
              ),

            if (widget.showUserData &&
                (widget.userGender != null ||
                    widget.userAge != null ||
                    widget.userGoal != null ||
                    widget.userHeightCm != null ||
                    widget.userWeightKg != null))
              const SizedBox(height: 16),

            // Error display
            if (widget.error != null)
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: theme.colorScheme.errorContainer,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Row(
                  children: [
                    Icon(
                      Icons.error_outline,
                      color: theme.colorScheme.onErrorContainer,
                    ),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        widget.error!,
                        style: TextStyle(
                          color: theme.colorScheme.onErrorContainer,
                        ),
                      ),
                    ),
                  ],
                ),
              ),

            if (widget.error != null) const SizedBox(height: 16),

            // Generate button
            SizedBox(
              width: double.infinity,
              child: ElevatedButton.icon(
                onPressed: widget.isLoading ? null : widget.onGenerate,
                icon: widget.isLoading
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Icon(Icons.restaurant_menu),
                label: Text(
                  widget.isLoading ? 'Generating...' : 'Generate Meal Plan',
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ProfileChip extends StatelessWidget {
  final String label;
  final String value;
  final IconData icon;

  const _ProfileChip({
    required this.label,
    required this.value,
    required this.icon,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: theme.colorScheme.primaryContainer.withValues(alpha: 0.3),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 14, color: theme.colorScheme.primary),
          const SizedBox(width: 4),
          Text(
            '$label: $value',
            style: theme.textTheme.labelSmall?.copyWith(
              color: theme.colorScheme.primary,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }
}

class WeeklyPlanGenerator extends StatefulWidget {
  final VoidCallback onGenerate;
  final bool isLoading;
  final String? error;
  final String? userGender;
  final int? userAge;
  final String? userGoal;
  final double? userHeightCm;
  final double? userWeightKg;
  final bool showUserData;

  const WeeklyPlanGenerator({
    super.key,
    required this.onGenerate,
    required this.isLoading,
    this.error,
    this.userGender,
    this.userAge,
    this.userGoal,
    this.userHeightCm,
    this.userWeightKg,
    this.showUserData = true,
  });

  @override
  State<WeeklyPlanGenerator> createState() => _WeeklyPlanGeneratorState();
}

class _WeeklyPlanGeneratorState extends State<WeeklyPlanGenerator> {
  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      margin: const EdgeInsets.symmetric(vertical: 8),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Weekly Meal Plan Generator',
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),

            // Show comprehensive user data if available
            if (widget.showUserData &&
                (widget.userGender != null ||
                    widget.userAge != null ||
                    widget.userGoal != null ||
                    widget.userHeightCm != null ||
                    widget.userWeightKg != null))
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: theme.colorScheme.primaryContainer.withValues(
                    alpha: 0.1,
                  ),
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(
                    color: theme.colorScheme.primary.withValues(alpha: 0.3),
                  ),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Icon(
                          Icons.person,
                          color: theme.colorScheme.primary,
                          size: 20,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          'Using your profile data:',
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color: theme.colorScheme.primary,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Wrap(
                      spacing: 16,
                      runSpacing: 8,
                      children: [
                        if (widget.userGender != null)
                          _ProfileChip(
                            label: 'Gender',
                            value: widget.userGender!,
                            icon: Icons.person_outline,
                          ),
                        if (widget.userAge != null)
                          _ProfileChip(
                            label: 'Age',
                            value: '${widget.userAge} years',
                            icon: Icons.cake,
                          ),
                        if (widget.userGoal != null)
                          _ProfileChip(
                            label: 'Goal',
                            value: widget.userGoal!
                                .replaceAll('_', ' ')
                                .toUpperCase(),
                            icon: Icons.flag,
                          ),
                        if (widget.userHeightCm != null)
                          _ProfileChip(
                            label: 'Height',
                            value:
                                '${widget.userHeightCm!.toStringAsFixed(0)} cm',
                            icon: Icons.height,
                          ),
                        if (widget.userWeightKg != null)
                          _ProfileChip(
                            label: 'Weight',
                            value:
                                '${widget.userWeightKg!.toStringAsFixed(1)} kg',
                            icon: Icons.monitor_weight,
                          ),
                      ],
                    ),
                  ],
                ),
              ),

            if (widget.showUserData &&
                (widget.userGender != null ||
                    widget.userAge != null ||
                    widget.userGoal != null ||
                    widget.userHeightCm != null ||
                    widget.userWeightKg != null))
              const SizedBox(height: 16),

            // Error display
            if (widget.error != null)
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: theme.colorScheme.errorContainer,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Row(
                  children: [
                    Icon(
                      Icons.error_outline,
                      color: theme.colorScheme.onErrorContainer,
                    ),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        widget.error!,
                        style: TextStyle(
                          color: theme.colorScheme.onErrorContainer,
                        ),
                      ),
                    ),
                  ],
                ),
              ),

            if (widget.error != null) const SizedBox(height: 16),

            // Generate button
            SizedBox(
              width: double.infinity,
              child: ElevatedButton.icon(
                onPressed: widget.isLoading ? null : widget.onGenerate,
                icon: widget.isLoading
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Icon(Icons.calendar_view_week),
                label: Text(
                  widget.isLoading ? 'Generating...' : 'Generate Weekly Plan',
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

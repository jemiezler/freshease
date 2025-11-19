import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../data/repositories/plans_repository.dart';
import '../../data/models/bundle_dtos.dart';

class PlansPage extends StatefulWidget {
  const PlansPage({super.key});

  @override
  State<PlansPage> createState() => _PlansPageState();
}

class _PlansPageState extends State<PlansPage> {
  final PlansRepository _repository = PlansRepository();
  List<BundleDTO> _bundles = [];
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadBundles();
  }

  Future<void> _loadBundles() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final bundles = await _repository.getActiveBundles();
      setState(() {
        _bundles = bundles;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  Color _getColorForIndex(int index) {
    final colors = [Colors.green, Colors.orange, Colors.teal, Colors.blue, Colors.purple];
    return colors[index % colors.length];
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(
          'FreshEase Plans',
          style: TextStyle(color: Colors.white),
        ),
        backgroundColor: Theme.of(context).colorScheme.primary,
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _error != null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        'Error loading plans',
                        style: TextStyle(color: Colors.red[700]),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        _error!,
                        style: TextStyle(color: Colors.grey[600]),
                        textAlign: TextAlign.center,
                      ),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: _loadBundles,
                        child: const Text('Retry'),
                      ),
                    ],
                  ),
                )
              : _bundles.isEmpty
                  ? const Center(
                      child: Text('No plans available'),
                    )
                  : LayoutBuilder(
                      builder: (context, constraints) {
                        final isWide = constraints.maxWidth >= 900;
                        final crossAxisCount = isWide
                            ? 3
                            : (constraints.maxWidth > 600 ? 2 : 1);
                        return GridView.builder(
                          padding: const EdgeInsets.all(16),
                          gridDelegate:
                              SliverGridDelegateWithFixedCrossAxisCount(
                            crossAxisCount: crossAxisCount,
                            crossAxisSpacing: 16,
                            mainAxisSpacing: 16,
                            childAspectRatio: isWide ? 1.1 : 0.95,
                          ),
                          itemCount: _bundles.length,
                          itemBuilder: (_, i) {
                            final bundle = _bundles[i];
                            final color = _getColorForIndex(i);
                            return _PlanCard(
                              bundle: bundle,
                              color: color,
                              onSubscribe: () => context.go(
                                '/plans/${bundle.id}',
                                extra: bundle,
                              ),
                            );
                          },
                        );
                      },
                    ),
    );
  }
}

class _PlanCard extends StatelessWidget {
  final BundleDTO bundle;
  final Color color;
  final VoidCallback onSubscribe;

  const _PlanCard({
    required this.bundle,
    required this.color,
    required this.onSubscribe,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 2,
      clipBehavior: Clip.antiAlias,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              decoration: BoxDecoration(
                color: color.withValues(alpha: .15),
                borderRadius: BorderRadius.circular(12),
              ),
              padding: const EdgeInsets.all(8),
              child: Icon(Icons.eco, color: color, size: 28),
            ),
            const SizedBox(height: 12),
            Text(
              bundle.name,
              style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w800),
            ),
            if (bundle.description != null && bundle.description!.isNotEmpty) ...[
              const SizedBox(height: 4),
              Text(
                bundle.description!,
                style: TextStyle(color: Colors.grey[600]),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
            ],
            const SizedBox(height: 12),
            Text(
              '฿${bundle.price.toStringAsFixed(0)} / plan',
              style: const TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.w700,
                color: Colors.green,
              ),
            ),
            const Spacer(),
            const SizedBox(height: 8),
            FilledButton.icon(
              onPressed: onSubscribe,
              icon: const Icon(Icons.arrow_forward),
              label: const Text('View Details'),
              style: FilledButton.styleFrom(
                minimumSize: const Size.fromHeight(44),
                backgroundColor: color.withValues(alpha: 0.9),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

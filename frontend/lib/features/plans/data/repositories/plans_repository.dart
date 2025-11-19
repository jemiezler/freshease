import 'package:get_it/get_it.dart';
import 'package:frontend/core/api/bundles_api.dart';
import 'package:frontend/core/api/bundle_items_api.dart';
import '../models/bundle_dtos.dart';

class PlansRepository {
  final BundlesApi _bundlesApi;
  final BundleItemsApi _bundleItemsApi;

  PlansRepository({
    BundlesApi? bundlesApi,
    BundleItemsApi? bundleItemsApi,
  })  : _bundlesApi = bundlesApi ?? GetIt.instance<BundlesApi>(),
        _bundleItemsApi = bundleItemsApi ?? GetIt.instance<BundleItemsApi>();

  /// Get all active bundles
  Future<List<BundleDTO>> getActiveBundles() async {
    try {
      final bundles = await _bundlesApi.listBundles();
      return bundles
          .map((json) => BundleDTO.fromJson(json))
          .where((bundle) => bundle.isActive)
          .toList();
    } catch (e) {
      throw Exception('Failed to fetch bundles: $e');
    }
  }

  /// Get bundle by ID
  Future<BundleDTO> getBundle(String bundleId) async {
    try {
      final bundle = await _bundlesApi.getBundle(bundleId);
      return BundleDTO.fromJson(bundle);
    } catch (e) {
      throw Exception('Failed to fetch bundle: $e');
    }
  }

  /// Get all bundle items
  Future<List<BundleItemDTO>> getAllBundleItems() async {
    try {
      final items = await _bundleItemsApi.listBundleItems();
      return items.map((json) => BundleItemDTO.fromJson(json)).toList();
    } catch (e) {
      throw Exception('Failed to fetch bundle items: $e');
    }
  }

  /// Get bundle items for a specific bundle
  Future<List<BundleItemDTO>> getBundleItems(String bundleId) async {
    try {
      final allItems = await getAllBundleItems();
      return allItems.where((item) => item.bundleId == bundleId).toList();
    } catch (e) {
      throw Exception('Failed to fetch bundle items: $e');
    }
  }

  /// Get bundle with its items
  Future<BundleWithItemsDTO> getBundleWithItems(String bundleId) async {
    try {
      final bundle = await getBundle(bundleId);
      final items = await getBundleItems(bundleId);
      return BundleWithItemsDTO(bundle: bundle, items: items);
    } catch (e) {
      throw Exception('Failed to fetch bundle with items: $e');
    }
  }
}


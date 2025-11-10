// core/api/api_di.dart
import 'package:get_it/get_it.dart';
import '../network/dio_client.dart';
import 'addresses_api.dart';
import 'bundles_api.dart';
import 'categories_api.dart';
import 'deliveries_api.dart';
import 'inventories_api.dart';
import 'meal_plans_api.dart';
import 'notifications_api.dart';
import 'orders_api.dart';
import 'payments_api.dart';
import 'products_api.dart';
import 'recipes_api.dart';
import 'reviews_api.dart';
import 'uploads_api.dart';
import 'vendors_api.dart';

/// Register all API clients in dependency injection
void registerApiClients(GetIt getIt) {
  // Core API clients
  getIt.registerLazySingleton<AddressesApi>(() => AddressesApi(getIt<DioClient>()));
  getIt.registerLazySingleton<BundlesApi>(() => BundlesApi(getIt<DioClient>()));
  getIt.registerLazySingleton<CategoriesApi>(() => CategoriesApi(getIt<DioClient>()));
  getIt.registerLazySingleton<DeliveriesApi>(() => DeliveriesApi(getIt<DioClient>()));
  getIt.registerLazySingleton<InventoriesApi>(() => InventoriesApi(getIt<DioClient>()));
  getIt.registerLazySingleton<MealPlansApi>(() => MealPlansApi(getIt<DioClient>()));
  getIt.registerLazySingleton<NotificationsApi>(() => NotificationsApi(getIt<DioClient>()));
  getIt.registerLazySingleton<OrdersApi>(() => OrdersApi(getIt<DioClient>()));
  getIt.registerLazySingleton<PaymentsApi>(() => PaymentsApi(getIt<DioClient>()));
  getIt.registerLazySingleton<ProductsApi>(() => ProductsApi(getIt<DioClient>()));
  getIt.registerLazySingleton<RecipesApi>(() => RecipesApi(getIt<DioClient>()));
  getIt.registerLazySingleton<ReviewsApi>(() => ReviewsApi(getIt<DioClient>()));
  getIt.registerLazySingleton<UploadsApi>(() => UploadsApi(getIt<DioClient>()));
  getIt.registerLazySingleton<VendorsApi>(() => VendorsApi(getIt<DioClient>()));
}


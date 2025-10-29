import 'package:get_it/get_it.dart';
import 'package:frontend/core/network/dio_client.dart';
import 'data/product_repository.dart';
import 'data/sources/shop_api.dart';

void registerShopDependencies(GetIt getIt) {
  // Register Shop API Service
  getIt.registerLazySingleton<ShopApiService>(
    () => ShopApiService(getIt<DioClient>().dio),
  );

  // Register Product Repository
  // You can switch between ApiProductRepository and MockProductRepository
  // For development, you might want to use MockProductRepository initially
  getIt.registerLazySingleton<ProductRepository>(
    () => ApiProductRepository(getIt<ShopApiService>()),
    // () => MockProductRepository(), // Uncomment this line to use mock data
  );
}

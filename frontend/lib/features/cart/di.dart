import 'package:get_it/get_it.dart';
import 'package:frontend/core/network/dio_client.dart';
import 'data/cart_repository.dart';
import 'data/sources/cart_api.dart';

void registerCartDependencies(GetIt getIt) {
  // Register Cart API Service
  getIt.registerLazySingleton<CartApiService>(
    () => CartApiService(getIt<DioClient>().dio),
  );

  // Register Cart Repository
  // You can switch between ApiCartRepository and MockCartRepository
  // For development, you might want to use MockCartRepository initially
  getIt.registerLazySingleton<CartRepository>(
    () => ApiCartRepository(getIt<CartApiService>()),
    // () => MockCartRepository(), // Uncomment this line to use mock data
  );
}

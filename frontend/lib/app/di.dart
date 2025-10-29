import 'package:frontend/core/health/health_controller.dart';
import 'package:frontend/core/health/health_repository.dart';
import 'package:frontend/core/health/health_service.dart';
import 'package:frontend/core/genai/genai_service.dart';
import 'package:frontend/core/state/cart_controller.dart';
import 'package:frontend/features/account/domain/repositories/user_repository.dart';
import 'package:frontend/features/account/presentation/state/user_cubit.dart';
import 'package:frontend/features/cart/data/cart_repository.dart';
import 'package:get_it/get_it.dart';
import '../core/network/dio_client.dart';
import '../features/auth/di.dart' as auth_di;
import '../features/account/di.dart' as account_di;
import '../features/shop/di.dart' as shop_di;
import '../features/cart/di.dart' as cart_di;
import 'package:flutter_dotenv/flutter_dotenv.dart';

final getIt = GetIt.instance;

Future<void> configureDependencies({String envFile = ".env"}) async {
  await dotenv.load(fileName: envFile);
  getIt.registerLazySingleton<DioClient>(() => DioClient());
  auth_di.registerAuthDependencies(getIt);
  account_di.registerAccountDependencies(getIt);
  shop_di.registerShopDependencies(getIt);
  cart_di.registerCartDependencies(getIt);
  if (!getIt.isRegistered<UserCubit>()) {
    getIt.registerLazySingleton<UserCubit>(
      () => UserCubit(getIt<UserRepository>()),
    );
  }
  getIt.registerSingleton<HealthService>(HealthService.instance);
  getIt.registerLazySingleton<GenAiService>(
    () => GenAiService(getIt<DioClient>()),
  );
  getIt.registerLazySingleton<HealthController>(
    () => HealthController(genAiService: getIt<GenAiService>()),
  );
  getIt.registerLazySingleton<HealthRepository>(() => NoopHealthRepository());

  // Register Cart Controller
  getIt.registerLazySingleton<CartController>(
    () => CartController(getIt<CartRepository>()),
  );

  await getIt<HealthController>().init();
}

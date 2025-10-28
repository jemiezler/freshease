import 'package:frontend/core/health/health_controller.dart';
import 'package:frontend/core/health/health_repository.dart';
import 'package:frontend/core/health/health_service.dart';
import 'package:get_it/get_it.dart';
import '../core/network/dio_client.dart';
import '../features/auth/di.dart' as auth_di;
import 'package:flutter_dotenv/flutter_dotenv.dart';

final getIt = GetIt.instance;

Future<void> configureDependencies({String envFile = ".env"}) async {
  await dotenv.load(fileName: envFile);
  getIt.registerLazySingleton<DioClient>(() => DioClient());
  auth_di.registerAuthDependencies(getIt);

  getIt.registerSingleton<HealthService>(HealthService.instance);
  getIt.registerLazySingleton<HealthController>(() => HealthController());
  getIt.registerLazySingleton<HealthRepository>(() => NoopHealthRepository());
  await getIt<HealthController>().init();
}

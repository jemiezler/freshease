import 'package:get_it/get_it.dart';
import '../core/network/dio_client.dart';
import '../features/auth/di.dart' as auth_di;

final getIt = GetIt.instance;

Future<void> configureDependencies() async {
  getIt.registerLazySingleton<DioClient>(() => DioClient());

  auth_di.registerAuthDependencies(getIt);
}

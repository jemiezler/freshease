import 'package:get_it/get_it.dart';
import '../../core/network/dio_client.dart';
import 'data/sources/auth_api.dart';
import 'data/repositories/auth_repository_impl.dart';
import 'domain/repositories/auth_repository.dart';

void registerAuthDependencies(GetIt getIt) {
  getIt
    ..registerLazySingleton<AuthApi>(() => AuthApi(getIt<DioClient>()))
    ..registerLazySingleton<AuthRepository>(
      () => AuthRepositoryImpl(getIt<AuthApi>()),
    );
}

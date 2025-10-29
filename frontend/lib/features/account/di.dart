// di.dart
import 'package:get_it/get_it.dart';
import '../../../core/network/dio_client.dart';
import 'data/sources/user_api.dart';
import 'data/repositories/user_repository_impl.dart';
import 'domain/repositories/user_repository.dart';
import 'presentation/state/user_cubit.dart';

void registerAccountDependencies(GetIt getIt) {
  // Register UserApi
  getIt.registerLazySingleton<UserApi>(() => UserApi(getIt<DioClient>()));

  // Register UserRepository
  getIt.registerLazySingleton<UserRepository>(
    () => UserRepositoryImpl(getIt<UserApi>()),
  );

  // Register UserCubit as singleton (same instance throughout app lifecycle)
  getIt.registerLazySingleton<UserCubit>(
    () => UserCubit(getIt<UserRepository>()),
  );
}

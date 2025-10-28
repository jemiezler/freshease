import 'package:get_it/get_it.dart';
import '../core/network/dio_client.dart';
import '../features/auth/di.dart' as auth_di;
import 'package:flutter_dotenv/flutter_dotenv.dart';

final getIt = GetIt.instance;

Future<void> configureDependencies({String envFile = ".env"}) async {
  await dotenv.load(fileName: envFile);
  getIt.registerLazySingleton<DioClient>(() => DioClient());
  auth_di.registerAuthDependencies(getIt);
}

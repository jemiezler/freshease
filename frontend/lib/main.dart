import 'package:flutter/material.dart';
import 'app/app.dart';
import 'app/di.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await configureDependencies(envFile: ".env");
  runApp(const App());
}

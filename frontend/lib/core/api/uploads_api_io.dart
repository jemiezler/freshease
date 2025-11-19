// mobile implementation
import 'dart:io';

Future<List<int>> readFileBytes(File file) async {
  return await file.readAsBytes();
}

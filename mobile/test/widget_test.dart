import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:qonaqzhai_mobile/main.dart';

void main() {
  testWidgets('App boots without throwing', (WidgetTester tester) async {
    await tester.runAsync(() async {
      await tester.pumpWidget(const ProviderScope(child: QonaqzhaiApp()));
      await tester.pump(const Duration(milliseconds: 50));
    });
    expect(find.byType(MaterialApp), findsOneWidget);
  });
}

class ChatMessage {
  ChatMessage({required this.role, required this.text, this.blocks = const []});

  final String role; // user | ai
  final String text;
  final List<Map<String, dynamic>> blocks;
}

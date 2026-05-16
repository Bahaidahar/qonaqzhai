class AppNotification {
  const AppNotification({
    required this.id,
    required this.type,
    required this.channel,
    required this.title,
    required this.body,
    required this.status,
    required this.createdAt,
  });

  final String id;
  final String type;
  final String channel;
  final String title;
  final String body;
  final String status;
  final String createdAt;
}

INSERT INTO users (email, name, hash) VALUES
    ('alice@example.com', 'Alice Johnson', '$2a$10$examplehash1234567890abcdefghijklmnopqrstuvwx'),
    ('bob@example.com',   'Bob Martinez',  '$2a$10$examplehash1234567890abcdefghijklmnopqrstuvwy'),
    ('carol@example.com', 'Carol Nguyen',  '$2a$10$examplehash1234567890abcdefghijklmnopqrstuvwz');

INSERT INTO videos (title, author_id, views, s3_key, size_bytes, duration_ms, status) VALUES
    ('My First Vlog',        1, 152,   'videos/11111111-1111-1111-1111-111111111111/original.mp4', 52428800,   305000, 'ready'),
    ('10 Minute Guitar Tips', 2, 8934, 'videos/22222222-2222-2222-2222-222222222222/original.mp4', 104857600, 612000, 'ready'),
    ('Cooking Pasta at Home', 3, 47,   'videos/33333333-3333-3333-3333-333333333333/original.mp4', 78643200,   480000, 'ready');

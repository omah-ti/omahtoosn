INSERT INTO tryouts (
    id,
    slug,
    title,
    description,
    instructions,
    status,
    duration_minutes,
    starts_at,
    ends_at,
    results_published_at,
    show_leaderboard
) VALUES (
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'to-osn-demo-v1',
    'TO-OSN Demo V1',
    'Demo tryout for local backend testing.',
    'Read each question carefully. Answers are saved on the server. Submit before the timer ends.',
    'ongoing',
    90,
    NOW() - INTERVAL '1 day',
    NOW() + INTERVAL '30 days',
    NOW(),
    TRUE
)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO questions (
    id,
    tryout_id,
    code,
    question_type,
    prompt_html,
    display_order,
    points,
    explanation_html
) VALUES
(
    '11111111-1111-1111-1111-111111111111',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'INF-001',
    'multiple_choice',
    '<p>What is the time complexity of binary search on a sorted array?</p>',
    1,
    10,
    '<p>Binary search cuts the search space in half at each step.</p>'
),
(
    '22222222-2222-2222-2222-222222222222',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'INF-002',
    'multiple_choice',
    '<p>Which data structure follows the LIFO principle?</p>',
    2,
    10,
    '<p>Stack uses Last In First Out.</p>'
),
(
    '33333333-3333-3333-3333-333333333333',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'INF-003',
    'short_text',
    '<p>How many bits are there in one byte?</p>',
    3,
    10,
    '<p>One byte equals eight bits.</p>'
)
ON CONFLICT (tryout_id, code) DO NOTHING;

INSERT INTO question_options (id, question_id, option_key, option_text, display_order) VALUES
('11111111-aaaa-aaaa-aaaa-aaaaaaaaaaa1', '11111111-1111-1111-1111-111111111111', 'A', 'O(1)', 1),
('11111111-aaaa-aaaa-aaaa-aaaaaaaaaaa2', '11111111-1111-1111-1111-111111111111', 'B', 'O(log n)', 2),
('11111111-aaaa-aaaa-aaaa-aaaaaaaaaaa3', '11111111-1111-1111-1111-111111111111', 'C', 'O(n)', 3),
('11111111-aaaa-aaaa-aaaa-aaaaaaaaaaa4', '11111111-1111-1111-1111-111111111111', 'D', 'O(n log n)', 4),
('11111111-aaaa-aaaa-aaaa-aaaaaaaaaaa5', '11111111-1111-1111-1111-111111111111', 'E', 'O(n^2)', 5),
('22222222-bbbb-bbbb-bbbb-bbbbbbbbbbb1', '22222222-2222-2222-2222-222222222222', 'A', 'Queue', 1),
('22222222-bbbb-bbbb-bbbb-bbbbbbbbbbb2', '22222222-2222-2222-2222-222222222222', 'B', 'Stack', 2),
('22222222-bbbb-bbbb-bbbb-bbbbbbbbbbb3', '22222222-2222-2222-2222-222222222222', 'C', 'Heap', 3),
('22222222-bbbb-bbbb-bbbb-bbbbbbbbbbb4', '22222222-2222-2222-2222-222222222222', 'D', 'Graph', 4),
('22222222-bbbb-bbbb-bbbb-bbbbbbbbbbb5', '22222222-2222-2222-2222-222222222222', 'E', 'Tree', 5)
ON CONFLICT (question_id, option_key) DO NOTHING;

INSERT INTO question_answer_keys (question_id, correct_option_key) VALUES
('11111111-1111-1111-1111-111111111111', 'B'),
('22222222-2222-2222-2222-222222222222', 'B')
ON CONFLICT (question_id) DO NOTHING;

INSERT INTO question_short_answer_variants (id, question_id, answer_text, normalized_text) VALUES
('33333333-cccc-cccc-cccc-ccccccccccc1', '33333333-3333-3333-3333-333333333333', '8', '8'),
('33333333-cccc-cccc-cccc-ccccccccccc2', '33333333-3333-3333-3333-333333333333', 'eight', 'eight')
ON CONFLICT (question_id, normalized_text) DO NOTHING;

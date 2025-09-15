-- DDD Bottom-Up サークル管理システム データベーススキーマ

-- ユーザーテーブル
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    is_premium BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_premium (is_premium)
);

-- サークルテーブル
CREATE TABLE circles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_id VARCHAR(36) NOT NULL,
    created_at DATETIME NOT NULL,
    member_count INT DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name),
    INDEX idx_owner_id (owner_id),
    INDEX idx_created_at (created_at),
    INDEX idx_member_count (member_count),
    INDEX idx_recommended (created_at, member_count), -- おすすめサークル用複合インデックス
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

-- サークルメンバーテーブル
CREATE TABLE circle_members (
    circle_id VARCHAR(36),
    user_id VARCHAR(36),
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (circle_id, user_id),
    INDEX idx_user_id (user_id),
    INDEX idx_joined_at (joined_at),
    FOREIGN KEY (circle_id) REFERENCES circles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- サークルメンバー数を自動更新するトリガー
DELIMITER $$

CREATE TRIGGER update_member_count_after_insert
AFTER INSERT ON circle_members
FOR EACH ROW
BEGIN
    UPDATE circles 
    SET member_count = (
        SELECT COUNT(*) FROM circle_members WHERE circle_id = NEW.circle_id
    )
    WHERE id = NEW.circle_id;
END$$

CREATE TRIGGER update_member_count_after_delete
AFTER DELETE ON circle_members
FOR EACH ROW
BEGIN
    UPDATE circles 
    SET member_count = (
        SELECT COUNT(*) FROM circle_members WHERE circle_id = OLD.circle_id
    )
    WHERE id = OLD.circle_id;
END$$

DELIMITER ;

-- 初期データ（テスト用）
INSERT INTO users (id, first_name, last_name, email, is_premium) VALUES
('user1', '太郎', '田中', 'taro@example.com', FALSE),
('user2', '花子', '鈴木', 'hanako@example.com', TRUE),
('user3', '次郎', '佐藤', 'jiro@example.com', TRUE);

INSERT INTO circles (id, name, owner_id, created_at, member_count) VALUES
('circle1', 'プログラミング勉強会', 'user1', DATE_SUB(NOW(), INTERVAL 15 DAY), 12),
('circle2', 'デザイン研究会', 'user2', DATE_SUB(NOW(), INTERVAL 45 DAY), 8),
('circle3', 'AI・機械学習サークル', 'user3', DATE_SUB(NOW(), INTERVAL 5 DAY), 15);
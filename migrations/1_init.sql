-- +migrate Up
CREATE TABLE `todos` (
  `todo_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `title` varchar(256) NOT NULL COMMENT 'タスクのタイトル',
  `comment` varchar(256) COMMENT 'タスクへのコメント',
  `completed` BOOL NOT NULL DEFAULT FALSE COMMENT 'タスクが完了したか否か',
  `created` datetime NOT NULL DEFAULT NOW() COMMENT '登録日',
  `updated` datetime DEFAULT NULL COMMENT '更新日',
  PRIMARY KEY (`todo_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='タスクリスト';

CREATE TABLE `users` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(256) NOT NULL COMMENT 'ユーザ名',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='ユーザリスト';

INSERT INTO users VALUES
  (1, 'hoge'),
  (2, 'fuga'),
  (3, 'moge');

-- +migrate Down
DROP TABLE todos;

# @TODO マイグレーションうまく実行できるように確認する

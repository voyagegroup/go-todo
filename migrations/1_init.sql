-- +migrate Up
CREATE TABLE `todos` (
  `todo_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `title` varchar(256) NOT NULL COMMENT 'タスクのタイトル',
  `completed` BOOL NOT NULL DEFAULT FALSE COMMENT 'タスクが完了したか否か',
  `created` datetime NOT NULL DEFAULT NOW() COMMENT '登録日',
  `updated` datetime DEFAULT NULL COMMENT '更新日',
  PRIMARY KEY (`todo_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='タスクリスト';

CREATE TABLE `todo_comments` (
  `comment_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `todo_id` int(11) NOT NULL COMMENT 'タスクID',
  `comment` varchar(256) NOT NULL COMMENT 'タスクへのコメント',
  `created` datetime NOT NULL DEFAULT NOW() COMMENT '登録日',
  `updated` datetime DEFAULT NULL COMMENT '更新日',
  PRIMARY KEY (`comment_id`),
  FOREIGN KEY (`todo_id`) REFERENCES `todos`(`todo_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='タスクコメントリスト';

CREATE TABLE `users` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(256) NOT NULL COMMENT 'ユーザー名',
  `created` datetime NOT NULL DEFAULT NOW() COMMENT '登録日',
  `updated` datetime DEFAULT NULL COMMENT '更新日',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='ユーザーリスト';

CREATE TABLE `reports` (
  `report_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `title` varchar(256) NOT NULL COMMENT 'タイトル',
  `body` text NOT NULL COMMENT '内容',
  `created` datetime NOT NULL DEFAULT NOW() COMMENT '登録日',
  PRIMARY KEY (`report_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='日報リスト';

CREATE TABLE `report_tagmaps` (
  `report_id` int(11) NOT NULL COMMENT '日報ID',
  `tag_id` int(11) NOT NULL COMMENT 'タグID',
  PRIMARY KEY (`report_id`, `tag_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='タグリスト';


CREATE TABLE `report_tags` (
  `tag_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(256) NOT NULL COMMENT 'タイトル',
  PRIMARY KEY (`tag_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='タグリスト';


-- +migrate Down
DROP TABLE todos;
DROP TABLE todo_comments;
DROP TABLE users;
DROP TABLE reports;
DROP TABLE report_tagmaps;
DROP TABLE report_tags;

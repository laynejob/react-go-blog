package db

import (
    "time"
    "fmt"
)

type JsonTime time.Time

func (t *JsonTime) MarshalJSON() ([]byte, error) {
    st := fmt.Sprintf(`"%s"`, time.Time(*t).Local().Format(timeFormart))
    return []byte(st), nil
}

type User struct {
    Id             uint64    `json:"id"`
    Username     string    `json:"username"`
    Password     string    `json:"-"`
    IsSuper     bool    `json:"-"`
    Nickname     string    `json:"nickname"`
    Avatar         string    `json:"avatar"`
    Email         string    `json:"email"`
    QQ             string    `json:"qq"`
    WeChat         string    `json:"wechat"`
    CTime         JsonTime    `json:"cTime"`
    LTime         JsonTime    `json:"lTime"`
}

var stmt = `
— —
— Table ‘topic’
—
— —

CREATE TABLE IF NOT EXISTS topic (
id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
name VARCHAR(255) NOT NULL,
avatar VARCHAR(255) NULL DEFAULT NULL,
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS category (
id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
name VARCHAR(255) NOT NULL,
parentId INT(11) UNSIGNED NULL DEFAULT NULL,
PRIMARY KEY (id),
KEY parentId (parentId),
CONSTRAINT category_parent_fk_1 FOREIGN KEY (parentId) REFERENCES category (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS page (
id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
userId INT(11) UNSIGNED NOT NULL,
topicId INT(11) UNSIGNED NULL DEFAULT NULL,
sequence INT(11) UNSIGNED NULL DEFAULT NULL,
md5 VARCHAR(32) NULL DEFAULT NULL,
title VARCHAR(255) NOT NULL,
summary VARCHAR(255) NOT NULL,
avatar VARCHAR(255) NULL DEFAULT NULL,
ctime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
mtime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
status TINYINT NULL DEFAULT NULL,
content TEXT NULL DEFAULT NULL,
draft TEXT NULL DEFAULT NULL,
readNum INT(11) NULL DEFAULT NULL,
likeNum INT(11) NULL DEFAULT NULL,
commentNum INT(11) NULL DEFAULT NULL,
PRIMARY KEY (id),
KEY userId (userId),
KEY topicId (topicId),
UNIQUE key md5 (md5),
CONSTRAINT page_user_fk_1 FOREIGN KEY (userId) REFERENCES user (id),
CONSTRAINT page_topic_fk_2 FOREIGN KEY (topicId) REFERENCES topic (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

— —
— Table ‘comment’
—
— —

CREATE TABLE IF NOT EXISTS comment (
id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
userId INT(11) UNSIGNED NOT NULL,
pageId INT(11) UNSIGNED NOT NULL,
parentId INT(11) UNSIGNED NULL DEFAULT NULL,
content VARCHAR(255) NOT NULL,
likeNum INT(11) NULL DEFAULT NULL,
commentNum INT(11) NULL DEFAULT NULL,
PRIMARY KEY (id),
KEY userId (userId),
KEY pageId (pageId),
KEY parentId (parentId),
CONSTRAINT comment_user_fk_1 FOREIGN KEY (userId) REFERENCES user (id),
CONSTRAINT comment_page_fk_2 FOREIGN KEY (pageId) REFERENCES page (id),
CONSTRAINT comment_parent_fk_2 FOREIGN KEY (parentId) REFERENCES comment (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS likepage (
id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
userId INT(11) UNSIGNED NOT NULL,
pageId INT(11) UNSIGNED NOT NULL,
PRIMARY KEY (id),
KEY userId (userId),
KEY pageId (pageId),
UNIQUE KEY userPageId (userId, pageId),
CONSTRAINT likepage_user_fk_1 FOREIGN KEY (userId) REFERENCES user (id),
CONSTRAINT likepage_page_fk_2 FOREIGN KEY (pageId) REFERENCES page (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS likecomment (
id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
userId INT(11) UNSIGNED NOT NULL,
commentId INT(11) UNSIGNED NOT NULL,
PRIMARY KEY (id),
KEY userId (userId),
KEY pageId (commentId),
UNIQUE KEY userPageId (userId, commentId),
CONSTRAINT likecomment_user_fk_1 FOREIGN KEY (userId) REFERENCES user (id),
CONSTRAINT likecomment_comment_fk_2 FOREIGN KEY (commentId) REFERENCES comment (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`
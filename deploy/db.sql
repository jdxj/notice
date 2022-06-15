CREATE TABLE `rss_urls`
(
    `id`  bigint       NOT NULL AUTO_INCREMENT,
    `url` varchar(500) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `rss_urls_UN` (`url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `github`
(
    `id`    bigint       NOT NULL AUTO_INCREMENT,
    `owner` varchar(100) NOT NULL,
    `repo`  varchar(100) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `owner_repo_idx` (`owner`,`repo`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

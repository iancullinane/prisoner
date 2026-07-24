# Changelog

## [1.4.0](https://github.com/iancullinane/prisoner/compare/v1.3.0...v1.4.0) (2026-07-10)


### Features

* add frontend container ([ebad2d9](https://github.com/iancullinane/prisoner/commit/ebad2d960bd000dc105731742f7435bd3c738a5e))
* add history API client ([1051199](https://github.com/iancullinane/prisoner/commit/1051199671350c132e19d267b6b4ea12a879a312))
* add HistoryPage ([283dfee](https://github.com/iancullinane/prisoner/commit/283dfee802c06c9a2fbe99827a21d94e7ae3fd0f))
* add Interaction/PlayRequest/PlayResponse frontend types ([e66bdf0](https://github.com/iancullinane/prisoner/commit/e66bdf03ef8349a9eac882c8c9524fa8738b7684))
* add play API client ([67a8768](https://github.com/iancullinane/prisoner/commit/67a87682b73d58362a44982d95f6b9b824c31a13))
* add PlayersPage ([adce0c2](https://github.com/iancullinane/prisoner/commit/adce0c2d8c355d3a0a82abb790ebf9864d149b70))
* add PlayPage ([6fcbbfe](https://github.com/iancullinane/prisoner/commit/6fcbbfe6b24a43b625ba01e05fe8232bfb14474c))
* add router, Tailwind theme, Sidebar, and Layout ([65fa8bd](https://github.com/iancullinane/prisoner/commit/65fa8bd8c32157e5e15ea82fbe3b28744d1f854b))
* point players API client at /api/v1 ([aad3f20](https://github.com/iancullinane/prisoner/commit/aad3f20259f1470f6b1cc23fd5eb70bd076c8e31))
* version backend API under /api/v1 ([06e6e46](https://github.com/iancullinane/prisoner/commit/06e6e465fe057aad13bd2a1181f1708ecf202981))
* wire up multi-page routing in App ([8185600](https://github.com/iancullinane/prisoner/commit/8185600cc7abd2f6ad2439b5c4065d5cbdff2d56))


### Bug Fixes

* add healtz endpoint ([4b19041](https://github.com/iancullinane/prisoner/commit/4b1904115448d54f299bf6c17d93e5b9fb1a38f5))
* style PlayPage move-selection buttons for visible selected state ([37d8140](https://github.com/iancullinane/prisoner/commit/37d814026efcad2d3ea0df8383721567e55dfdfc))
* update json casing to be consistent ([9b1f877](https://github.com/iancullinane/prisoner/commit/9b1f877c83a8058294d011a7f79eb2932445bc5c))

## [1.3.0](https://github.com/iancullinane/prisoner/compare/v1.2.1...v1.3.0) (2026-07-06)


### Features

* add web frontend ([bece50c](https://github.com/iancullinane/prisoner/commit/bece50c104da0e27e30fbfe894b0333fc9d54a5b))


### Bug Fixes

* add analysis page ([1450979](https://github.com/iancullinane/prisoner/commit/1450979bccab70b951deb8bcc1bf72503dba46f6))
* add pages deploy of ANALYSIS.html ([aa42a74](https://github.com/iancullinane/prisoner/commit/aa42a746c6b2365b2e3aeb0fa88d7238a429abf4))
* add run dev for holistic testing ([1f7e6e6](https://github.com/iancullinane/prisoner/commit/1f7e6e68ada5b5d013ae166a418383397b61bf79))
* build `sqlc` files in the container ([4c3c91a](https://github.com/iancullinane/prisoner/commit/4c3c91a5668862b40b45d6726ba1e726cc300be9))
* handle CORS on server side ([ce9731e](https://github.com/iancullinane/prisoner/commit/ce9731ec040a50a83c9ea5ac1517c4a41e4e4d67))
* publish to ecr on like everything ([b2f25cd](https://github.com/iancullinane/prisoner/commit/b2f25cd7daaaef20760fc03c1e208c1272816283))
* update ECR repository to adventurebrave ([d20990e](https://github.com/iancullinane/prisoner/commit/d20990e81a3303306837d476fa36990d41f85c67))

## [1.2.1](https://github.com/iancullinane/prisoner/compare/v1.2.0...v1.2.1) (2026-07-06)


### Bug Fixes

* Add Claude review to CI flow ([1af37c1](https://github.com/iancullinane/prisoner/commit/1af37c190b0c3415291aefa288cd22e4963828e0))
* add ecr push for real ([1af37c1](https://github.com/iancullinane/prisoner/commit/1af37c190b0c3415291aefa288cd22e4963828e0))

## [1.2.0](https://github.com/iancullinane/prisoner/compare/v1.1.0...v1.2.0) (2026-07-05)


### Features

* new command tree topology, wire up json flag ([cb5f158](https://github.com/iancullinane/prisoner/commit/cb5f158adda503202ce1b2a73db88288e4c1afbd))
* Update stores and API with CreadtedAt for player field ([6e6a53a](https://github.com/iancullinane/prisoner/commit/6e6a53a594aab947d4cb52e6ad6d85bf6f0b4276))


### Bug Fixes

* add the printInteracton method on interactions ([290e383](https://github.com/iancullinane/prisoner/commit/290e3830f1e71d9c2c3898c6b874c573abece1aa))
* better errors out of postgres store ([ce0a8aa](https://github.com/iancullinane/prisoner/commit/ce0a8aabe1ac790a571e7c3a52b1f3d7723d56c6))
* build-and-test builds sqlc ([08f42c1](https://github.com/iancullinane/prisoner/commit/08f42c10f5a303b0bb6f9540ab9a924ae15c58f8))
* exclude sqlc files from git ([b07ce1b](https://github.com/iancullinane/prisoner/commit/b07ce1b448280c6596c047fc36adff0200490f40))
* give claude code review the apprpriate permissions ([2b4f080](https://github.com/iancullinane/prisoner/commit/2b4f080f7a1751619af2d147be3ffd0b46f94344))
* log store type earlier and better on player ([9f472d0](https://github.com/iancullinane/prisoner/commit/9f472d018d16badad36fc2c15b8f8a81064a7d01))
* update postgres store with player:CreatedAt ([543e32f](https://github.com/iancullinane/prisoner/commit/543e32f8b1963fb5a8744f3306638aa73e76e630))

## [1.1.0](https://github.com/iancullinane/prisoner/compare/v1.0.0...v1.1.0) (2026-07-04)


### Features

* add migrate command, rm migrate from server startup ([5270c7c](https://github.com/iancullinane/prisoner/commit/5270c7c82a9f350f44734929a210bfebca4005c5))


### Bug Fixes

* add test to interal/types/interaction ([63e7eee](https://github.com/iancullinane/prisoner/commit/63e7eee8944b874ef6299c48ac33d6809b979de3))
* allow bot-triggered runs in claude-code-review workflow ([6998c7e](https://github.com/iancullinane/prisoner/commit/6998c7e76d0013916e896c39700cafd2ef60e268))
* grant write permissions for pull-requests and issues in claude-code-review workflow ([134a9b8](https://github.com/iancullinane/prisoner/commit/134a9b85b3549c8a377b50d4ec3d72a047525dab))
* improve test coverage in prisoner package ([cb22b9b](https://github.com/iancullinane/prisoner/commit/cb22b9bd80da9344b976f186872ec7e69056c121))
* migrate no longer runs on container start ([2d79f7b](https://github.com/iancullinane/prisoner/commit/2d79f7b201f88c0a3d2e46777cada52480ab474d))
* rename to PlayedAt ([245eaad](https://github.com/iancullinane/prisoner/commit/245eaad70b305c6469870022045ce52611af0f84))
* unblock claude-review CI — bot actor + write permissions ([cca129e](https://github.com/iancullinane/prisoner/commit/cca129e85fcf3585607d3f609adbf32cce218da2))
* validate Move and Result internally ([f3fab2f](https://github.com/iancullinane/prisoner/commit/f3fab2f0696ab0d82f158047d054d857971f7726))

## [1.1.0](https://github.com/iancullinane/prisoner/compare/v1.0.0...v1.1.0) (2026-07-04)


### Features

* add migrate command, rm migrate from server startup ([5270c7c](https://github.com/iancullinane/prisoner/commit/5270c7c82a9f350f44734929a210bfebca4005c5))


### Bug Fixes

* add test to interal/types/interaction ([63e7eee](https://github.com/iancullinane/prisoner/commit/63e7eee8944b874ef6299c48ac33d6809b979de3))
* improve test coverage in prisoner package ([cb22b9b](https://github.com/iancullinane/prisoner/commit/cb22b9bd80da9344b976f186872ec7e69056c121))
* migrate no longer runs on container start ([2d79f7b](https://github.com/iancullinane/prisoner/commit/2d79f7b201f88c0a3d2e46777cada52480ab474d))
* rename to PlayedAt ([245eaad](https://github.com/iancullinane/prisoner/commit/245eaad70b305c6469870022045ce52611af0f84))
* validate Move and Result internally ([f3fab2f](https://github.com/iancullinane/prisoner/commit/f3fab2f0696ab0d82f158047d054d857971f7726))

## 1.0.0 (2026-06-18)


### Features

* fully commit to logging ([dc7d86d](https://github.com/iancullinane/prisoner/commit/dc7d86df8b6fe80320506860e9b30b871167bd2d))


### Bug Fixes

* getOrCreatePlayer now has multiple error type resturns ([47d120c](https://github.com/iancullinane/prisoner/commit/47d120c35ce2a07677bd1b052e8347e81ab77820))
* move play logic into POST body ([f21141e](https://github.com/iancullinane/prisoner/commit/f21141e1aeb6143a2b783c0c6902d2cfb7386873))
* test player can not play themselves ([0eec9a7](https://github.com/iancullinane/prisoner/commit/0eec9a74d47861941563383c7610324e324e83c1))

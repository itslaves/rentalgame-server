# rentalgames-server > config

서비스 환경 별로 설정 파일을 분리해서 관리한다.

- `production`: 서비스 MySQL, Redis 와 연동된 API 서버 환경
- `staging`: 서비스 MySQL, Redis 와 연동되지만 격리된 API 서버 환경
- `test`: 테스트 MySQL, Redis 와 연동된 API 서버 환경
- `develop`: 로컬 MySQL, Redis 와 연동된 API 서버 환경
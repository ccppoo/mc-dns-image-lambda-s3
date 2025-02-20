# MC DNS Image Lambda

## 수정사항 있는 경우 빌드 먼저하기

```bash
sam build
```

## 로컬에서 테스트용으로 실행하기

```bash
sam local start-api
```

## 동시에

```bash
sam build;  sam local start-api
```

## CICD

`github actions`에서 컴파일 후 zip 형태로 `lambda` 업데이트하는 형식으로 진행됨

컴파일 된 후 런타임에서 `ENV` 불러오는 방식임

### AWS Lambda Config

`AWS Lambda Configuration`에서 `ENV` 불러오는 방식

### Local Dev

`template.yml` 정의된 양식으로 컴파일함

`Resources > DemoFunction > Properties > Environment > Variables > ALLOW_ORIGINS`에서 환경변수 불러오기

# Moapick

## Description

Moapick Service의 server입니다.

Go로 작성되었으며 Fiber Framework로 구동됩니다.

Database는 PostgreSQL입니다. GORM을 사용하고 있습니다.

API 문서는 BASE_PATH/docs에서 확인하실 수 있습니다.

## HELP

아래는 개발시 맞닥뜨릴 수 있는 오류 상황입니다.

### swag mac에서 동작 안할 때

export PATH=$(go env GOPATH)/bin:$PATH

### swag이 API그룹별 문서를 생성하지 않을 때

-d 옵션을 사용해 directory를 ,를 기준으로 모두 명시해주어야합니다. 다만, main.go이 있는 루트도 포함해야합니다.

또 한, 타입 정의에 dependency가 포함되어있을 경우 --parseDependency를, internal을 포함한다면 --parseInternal옵션을 붙여야합니다. 이 옵션들은 동시에 사용될 수도 있습니다.

ex) swag init -d ./,./user,./article --parseDependency --parseInternal

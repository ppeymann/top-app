.PHONY: swagger
swagger:
	swag init --parseDependency --parseInternal -g /server/server.go

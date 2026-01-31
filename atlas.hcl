data "external_schema" "gorm" {
  program = ["go", "run", "./scripts/migrate.go"]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = getenv("DB_DSN_SHADOW")

  migration {
    dir = "file://database/migrations?format=goose"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

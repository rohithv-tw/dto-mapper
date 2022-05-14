# dto-mapper

Uses gjson and reflect to map from one dto to another dto

How to use:
Create struct with exposed fields and with tag `map_from`
with tag value as gjson tag that can be used to parse from source struct.

Example:

Refer `dto/source_dto.go` and `dto/target_dto.go`
Refer `main.go` for sample implementation.

Pending:

Add support for `time.Time` fields.
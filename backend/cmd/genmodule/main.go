package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// -------- config you may tweak --------
const ModulePath = "freshease/backend" // your go module path (from go.mod)

// -------- helpers --------

type Data struct {
	Name        string // plural module name (e.g., "roles")
	Route       string // "/roles"
	CapPlural   string // "Roles"
	Singular    string // "role"
	CapSing     string // "Role"
	SingularPkg string // ent entity import pkg (lowercase singular) e.g., "role"
	Backtick    string // "`"
	ModulePath  string // module import path
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./cmd/genmodule <module-name-plural>")
		os.Exit(1)
	}

	plural := strings.ToLower(os.Args[1])
	singular := naiveSingular(plural)
	caser := cases.Title(language.English)

	d := Data{
		Name:        plural,
		Route:       "/" + plural,
		CapPlural:   caser.String(plural),
		Singular:    singular,
		CapSing:     caser.String(singular),
		SingularPkg: singular,
		Backtick:    "`",
		ModulePath:  ModulePath,
	}

	// 1) module files
	modDir := filepath.Join("modules", d.Name)
	must(os.MkdirAll(modDir, 0o755))
	createAll(modDir, map[string]string{
		"controller.go": controllerTmpl,
		"service.go":    serviceTmpl,
		"repository.go": repoIfaceTmpl,
		"repo_ent.go":   repoEntTmpl,
		"dto.go":        dtoTmpl,
		"routes.go":     routesTmpl,
		"module.go":     moduleTmpl,
	}, d)

	// 2) ent scaffolding (schema + generate.go) if missing
	must(os.MkdirAll(filepath.Join("ent", "schema"), 0o755))
	createIfMissing(filepath.Join("ent", "generate.go"), entGenerateTmpl, d)
	createIfMissing(filepath.Join("ent", "schema", d.Singular+".go"), entSchemaTmpl, d)

	fmt.Println("Module scaffolded at:", modDir)
	fmt.Println("ℹ️  Next steps:")
	fmt.Println("   1) go generate ./ent")
	fmt.Println("   2) app.RegisterRoutes(...):", d.Name+".RegisterModuleWithEnt(api, entClient)")
}

func naiveSingular(plural string) string {
	if strings.HasSuffix(plural, "ies") {
		return plural[:len(plural)-3] + "y"
	}
	if strings.HasSuffix(plural, "ses") || strings.HasSuffix(plural, "xes") {
		return plural[:len(plural)-2]
	}
	if strings.HasSuffix(plural, "s") && len(plural) > 1 {
		return plural[:len(plural)-1]
	}
	return plural
}

func createAll(dir string, files map[string]string, d Data) {
	for name, src := range files {
		path := filepath.Join(dir, name)
		createIfMissing(path, src, d)
	}
}

func createIfMissing(path, tmpl string, d Data) {
	if fileExists(path) {
		fmt.Println("skip (exists):", path)
		return
	}
	if err := renderToFile(tmpl, d, path); err != nil {
		panic(err)
	}
	fmt.Println("created:", path)
}

func renderToFile(tmpl string, data any, path string) error {
	t := template.Must(template.New("f").Parse(tmpl))
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, data)
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// ---------------- templates ----------------

const controllerTmpl = `package {{.Name}}

import (
	"{{.ModulePath}}/internal/common/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/",   ctl.List{{.CapPlural}})
	r.Get("/:id", ctl.Get{{.CapSing}})
	r.Post("/",  ctl.Create{{.CapSing}})
	r.Patch("/:id", ctl.Update{{.CapSing}})
	r.Delete("/:id", ctl.Delete{{.CapSing}})
}

func (ctl *Controller) List{{.CapPlural}}(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "{{.CapPlural}} Retrieved Successfully"})
}

func (ctl *Controller) Get{{.CapSing}}(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "{{.CapSing}} Retrieved Successfully"})
}

func (ctl *Controller) Create{{.CapSing}}(c *fiber.Ctx) error {
	var dto Create{{.CapSing}}DTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "{{.CapSing}} Created Successfully"})
}

func (ctl *Controller) Update{{.CapSing}}(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto Update{{.CapSing}}DTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "{{.CapSing}} Updated Successfully"})
}

func (ctl *Controller) Delete{{.CapSing}}(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "{{.CapSing}} Deleted Successfully"})
}
`

const serviceTmpl = `package {{.Name}}

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context) ([]*Get{{.CapSing}}DTO, error)
	Get(ctx context.Context, id uuid.UUID) (*Get{{.CapSing}}DTO, error)
	Create(ctx context.Context, dto Create{{.CapSing}}DTO) (*Get{{.CapSing}}DTO, error)
	Update(ctx context.Context, id uuid.UUID, dto Update{{.CapSing}}DTO) (*Get{{.CapSing}}DTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service { return &service{repo: r} }

func (s *service) List(ctx context.Context) ([]*Get{{.CapSing}}DTO, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*Get{{.CapSing}}DTO, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, dto Create{{.CapSing}}DTO) (*Get{{.CapSing}}DTO, error) {
	return s.repo.Create(ctx, &dto)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, dto Update{{.CapSing}}DTO) (*Get{{.CapSing}}DTO, error) {
	dto.ID = id
	return s.repo.Update(ctx, &dto)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
`

const repoIfaceTmpl = `package {{.Name}}

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	List(ctx context.Context) ([]*Get{{.CapSing}}DTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Get{{.CapSing}}DTO, error)
	Create(ctx context.Context, u *Create{{.CapSing}}DTO) (*Get{{.CapSing}}DTO, error)
	Update(ctx context.Context, u *Update{{.CapSing}}DTO) (*Get{{.CapSing}}DTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
`

const repoEntTmpl = `package {{.Name}}

import (
	"context"

	"{{.ModulePath}}/ent"
	"{{.ModulePath}}/ent/{{.SingularPkg}}"
	"{{.ModulePath}}/internal/common/errs"
	"github.com/google/uuid"
)

type EntRepo struct{ c *ent.Client }

func NewEntRepo(client *ent.Client) Repository { return &EntRepo{c: client} }

func (r *EntRepo) List(ctx context.Context) ([]*Get{{.CapSing}}DTO, error) {
	rows, err := r.c.{{.CapSing}}.Query().Order(ent.Asc({{.SingularPkg}}.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*Get{{.CapSing}}DTO, 0, len(rows))
	for _, v := range rows {
		out = append(out, &Get{{.CapSing}}DTO{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
		})
	}
	return out, nil
}

func (r *EntRepo) FindByID(ctx context.Context, id uuid.UUID) (*Get{{.CapSing}}DTO, error) {
	v, err := r.c.{{.CapSing}}.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Get{{.CapSing}}DTO{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description,
	}, nil
}

func (r *EntRepo) Create(ctx context.Context, dto *Create{{.CapSing}}DTO) (*Get{{.CapSing}}DTO, error) {
	q := r.c.{{.CapSing}}.
		Create().
		SetName(dto.Name).
		SetDescription(dto.Description)

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &Get{{.CapSing}}DTO{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, dto *Update{{.CapSing}}DTO) (*Get{{.CapSing}}DTO, error) {
	q := r.c.{{.CapSing}}.UpdateOneID(dto.ID)

	if dto.Name != nil {
		q.SetName(*dto.Name)
	}
	if dto.Description != nil {
		q.SetDescription(*dto.Description)
	}

	if len(q.Mutation().Fields()) == 0 {
		return nil, errs.NoFieldsToUpdate
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &Get{{.CapSing}}DTO{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
	}, nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.c.{{.CapSing}}.DeleteOneID(id).Exec(ctx)
}
`

const dtoTmpl = `package {{.Name}}

import "github.com/google/uuid"

type Create{{.CapSing}}DTO struct {
	ID          uuid.UUID {{.Backtick}}json:"id" validate:"required"{{.Backtick}}
	Name        string    {{.Backtick}}json:"name" validate:"required,min=2,max=60"{{.Backtick}}
	Description string    {{.Backtick}}json:"description" validate:"required"{{.Backtick}}
}

type Update{{.CapSing}}DTO struct {
	ID          uuid.UUID {{.Backtick}}json:"id" validate:"required"{{.Backtick}}
	Name        *string   {{.Backtick}}json:"name" validate:"omitempty,min=2,max=60"{{.Backtick}}
	Description *string   {{.Backtick}}json:"description" validate:"omitempty"{{.Backtick}}
}

type Get{{.CapSing}}DTO struct {
	ID          uuid.UUID {{.Backtick}}json:"id" validate:"required"{{.Backtick}}
	Name        string    {{.Backtick}}json:"name" validate:"required,min=2,max=60"{{.Backtick}}
	Description string    {{.Backtick}}json:"description" validate:"required"{{.Backtick}}
}
`

const routesTmpl = `package {{.Name}}

import "github.com/gofiber/fiber/v2"

// Routes keeps routes isolated from wiring; controller methods attach here.
func Routes(app fiber.Router, ctl *Controller) {
	grp := app.Group("{{.Route}}")
	ctl.Register(grp)
}
`

const moduleTmpl = `package {{.Name}}

import (
	"github.com/gofiber/fiber/v2"
	"{{.ModulePath}}/ent"
)

// RegisterModuleWithEnt wires Ent repo -> service -> controller and mounts routes.
func RegisterModuleWithEnt(api fiber.Router, client *ent.Client) {
	repo := NewEntRepo(client)
	svc  := NewService(repo)
	ctl  := NewController(svc)
	Routes(api, ctl)
}
`

// ---- ent scaffolding ----

const entGenerateTmpl = `package ent

//go:generate go run entgo.io/ent/cmd/ent generate ./schema
`

const entSchemaTmpl = `package schema

import (
	"github.com/google/uuid"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type {{.CapSing}} struct{ ent.Schema }

func ({{.CapSing}}) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.String("name").NotEmpty(),
		field.String("description").NotEmpty(),
	}
}
`

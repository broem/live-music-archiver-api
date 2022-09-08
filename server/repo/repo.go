package repo

import (
	"encoding/json"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
)

type Repository struct {
	db *pg.DB
}

type Config struct {
	Address  string
	User     string
	Pass     string
	Database string
	Pool     int
}

func NewRepo(cfg Config) (*Repository, error) {
	db := pg.Connect(&pg.Options{
		Addr:     cfg.Address,
		User:     cfg.User,
		Password: cfg.Pass,
		Database: cfg.Database,
		// TLSConfig:             &tls.Config{},
		PoolSize: cfg.Pool,
	})

	return &Repository{
		db: db,
	}, nil

}

// AddUser adds a user to the database if they do not already exist
func (r *Repository) AddUser(user *User) error {
	if user.Email == "" {
		// idk how do you handle this?
		return errors.New("email is required")
	}
	count, err := r.db.Model(user).Where("email = ?", user.Email).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	_, err = r.db.Model(user).Insert()
	if err != nil {
		return err
	}

	return nil
}

// UpsertRunner upserts a runner to the database
func (r *Repository) UpsertRunner(runner *Runner) error {
	_, err := r.db.Model(runner).OnConflict("(user_id, map_id) DO UPDATE").Insert()
	if err != nil {
		return err
	}

	return nil
}

// UpdateRunner updates a runner by user_id and map_id with enabled
func (r *Repository) UpdateRunner(userID string, mapID uuid.UUID, enabled bool) error {
	_, err := r.db.Model(&Runner{
		UserID: userID,
		MapID:  mapID,
	}).Where("user_id = ? AND map_id = ?", userID, mapID).Set("enabled = ?", enabled).Update()
	if err != nil {
		return err
	}

	return nil
}

// DeleteRunner deletes a runner by user_id and map_id
func (r *Repository) DeleteRunner(userID, mapID string) error {
	_, err := r.db.Model(&Runner{
		UserID: userID,
		MapID:  uuid.MustParse(mapID),
	}).Where("user_id = ? AND map_id = ?", userID, mapID).Delete()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetRunnersByID(userID string) ([]*Runner, error) {
	var runners []*Runner
	err := r.db.Model(&runners).Where("user_id = ?", userID).Select()
	if err != nil {
		return nil, err
	}

	return runners, nil
}

// GetRunners returns all enabled runners
func (r *Repository) GetRunners() ([]*Runner, error) {
	var runners []*Runner
	err := r.db.Model(&runners).Where("enabled = true").Select()
	if err != nil {
		return nil, err
	}

	return runners, nil
}

// DisableRunner disables a runner by user_id and map_id
func (r *Repository) DisableRunner(userID, mapID string) error {
	_, err := r.db.Model(&Runner{
		UserID: userID,
		MapID:  uuid.MustParse(mapID),
	}).Where("user_id = ? AND map_id = ?", userID, mapID).Update(map[string]interface{}{
		"enabled": false,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) SaveScrapeBuilder(s *ScrapeBuilder) error {

	// turn scrape builder to json
	bulk, err := json.Marshal(s)
	if err != nil {
		return err
	}

	mapped := &Builder{
		UserID:        s.UserEmail,
		BuilderMap:   string(bulk),
	}

	_, err = r.db.Model(mapped).Insert()
	if err != nil {
		return err
	}

	return nil
}

// EventMapperExists checks if a event mapper exists using userId and venueBaseURL
func (r *Repository) EventMapperExists(userID, venueBaseURL string) (bool, error) {
	count, err := r.db.Model(&EventMapper{}).Where("user_id = ? AND venue_base_url = ?", userID, venueBaseURL).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// UpsertEventMapper upserts a event mapper to the database
func (r *Repository) UpsertEventMapper(m *EventMapper) (*EventMapper, error) {
	mapper := new(EventMapper)
	// check if event mapper exists
	exists, err := r.EventMapperExists(m.UserID, m.VenueBaseURL)
	if err != nil {
		return nil, err
	}

	if exists {
		_, err = r.db.Model(m).Where("user_id = ? AND venue_base_url = ?", m.UserID, m.VenueBaseURL).Update()
		if err != nil {
			return nil, err
		}

		// select the updated mapper
		err = r.db.Model(mapper).Where("user_id = ? AND venue_base_url = ?", m.UserID, m.VenueBaseURL).Select()
		if err != nil {
			return nil, err
		}

		return mapper, nil
	} else {
		_, err = r.db.Model(m).Insert()
		if err != nil {
			return nil, err
		}

		return m, nil
	}
}


func (r *Repository) SaveEventMapper(m *EventMapper) error {
	// assure we have a map ID
	if len(m.MapID) == 0 {
		m.MapID = uuid.Must(uuid.NewRandom())
	}

	_, err := r.db.Model(m).Insert()
	if err != nil {
		return err
	}

	return nil
}

// DeleteEventMapper deletes a event mapper by map_id
func (r *Repository) DeleteEventMapper(mapID string) error {
	_, err := r.db.Model(&EventMapper{
		MapID: uuid.MustParse(mapID),
	}).Where("map_id = ?", mapID).Delete()
	if err != nil {
		return err
	}

	return nil
}

// UpdateEventMapper updates a event mapper by map_id to approved
func (r *Repository) UpdateEventMapper(mapID uuid.UUID, enabled bool) error {
	_, err := r.db.Model(&EventMapper{
		MapID: mapID,
	}).Where("map_id = ?", mapID).Set("approved = ?", enabled).Update()
	if err != nil {
		return err
	}

	return nil
}

// SaveEvents saves a list of events to the database
func (r *Repository) SaveEvents(events []*Event) error {
	for _, e := range events {
		_, err := r.db.Model(e).Insert()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetEventsByID returns all events for a given user_id
func (r *Repository) GetEventsByID(userID string) ([]*Event, error) {
	var events []*Event
	err := r.db.Model(&events).Where("user_id = ?", userID).Select()
	if err != nil {
		return nil, err
	}

	return events, nil
}

// DeleteEvents deletes events by map_id
func (r *Repository) DeleteEvents(mapID string) error {
	_, err := r.db.Model(&Event{
		MapID: uuid.MustParse(mapID),
	}).Where("map_id = ?", mapID).Delete()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetScrapeBuilders(userID string) ([]*ScrapeBuilder, error) {
	var builders []*ScrapeBuilder
	err := r.db.Model(&builders).Where("user_id = ?", userID).Select()
	if err != nil {
		return nil, err
	}

	return builders, nil
}

// GetEventMapper returns the EventMapper for a given user_id and venue_base_url
func (r *Repository) GetEventMapper(userID, venueBaseURL string) (*EventMapper, error) {
	var mapper EventMapper
	err := r.db.Model(&mapper).Where("user_id = ? AND venue_base_url = ?", userID, venueBaseURL).Select()
	if err != nil {
		return nil, err
	}

	return &mapper, nil
}

// GetEventMappers returns all EventMappers for a given user_id
func (r *Repository) GetEventMappers(userID string) ([]*EventMapper, error) {
	var mappers []*EventMapper
	err := r.db.Model(&mappers).Where("user_id = ? AND approved = true", userID).Select()
	if err != nil {
		return nil, err
	}

	return mappers, nil
}

// GetEventMapperByID returns the EventMapper for a given map_id
func (r *Repository) GetEventMapperByID(mapID string) (*EventMapper, error) {
	var mapper EventMapper
	err := r.db.Model(&mapper).Where("map_id = ? AND approved = true", mapID).Select()
	if err != nil {
		return nil, err
	}

	return &mapper, nil
}

// GetEventMapperByMapID returns collection of EventMapper for a given map_ids
func (r *Repository) GetEventMappersByMapID(mapIDs []string) ([]*EventMapper, error) {
	var mappers []*EventMapper
	err := r.db.Model(&mappers).Where("map_id IN (?)", pg.In(mapIDs)).Select()
	if err != nil {
		return nil, err
	}

	return mappers, nil
}

// GetIgCapturedByUserID returns all IgCaptured for a given user_id
func (r *Repository) GetIgCapturedByUserID(userID string) ([]*IgCaptured, error) {
	var captures []*IgCaptured
	err := r.db.Model(&captures).Where("user_id = ?", userID).Select()
	if err != nil {
		return nil, err
	}

	return captures, nil
}

// SaveIgMap saves an IgMapper to the database
func (r *Repository) SaveIgMapper(m *IgMapper) error {
	// assure we have a map ID
	if len(m.MapID) == 0 {
		m.MapID = uuid.Must(uuid.NewRandom())
	}

	_, err := r.db.Model(m).Insert()
	if err != nil {
		return err
	}

	return nil
}

// GetIgMapperByID returns the IgMapper for a given map_id
func (r *Repository) GetIgMapperByID(mapID string) (*IgMapper, error) {
	var mapper IgMapper
	err := r.db.Model(&mapper).Where("map_id = ?", mapID).Select()
	if err != nil {
		return nil, err
	}

	return &mapper, nil
}

// UpsertIgRunner upserts an IgRunner to the database
func (r *Repository) UpsertIgRunner(runner *IgRunner) error {
	_, err := r.db.Model(runner).OnConflict("(user_id, map_id) DO UPDATE").Insert()
	if err != nil {
		return err
	}

	return nil
}

// GetIgRunners returns all enabled runners
func (r *Repository) GetIgRunners() ([]*IgRunner, error) {
	var runners []*IgRunner
	err := r.db.Model(&runners).Where("enabled = true").Select()
	if err != nil {
		return nil, err
	}

	return runners, nil
}

// GetIgRunnersByMapID returns all enabled runners for a given map_id
func (r *Repository) GetIgRunnersByUserID(userID string) ([]*IgRunner, error) {
	var runners []*IgRunner
	err := r.db.Model(&runners).Where("user_id = ?", userID).Select()
	if err != nil {
		return nil, err
	}

	return runners, nil
}

func (r *Repository) GetIGEventMappersByMapID(mapIDs []string) ([]*IgMapper, error) {
	var mappers []*IgMapper
	err := r.db.Model(&mappers).Where("map_id IN (?)", pg.In(mapIDs)).Select()
	if err != nil {
		return nil, err
	}

	return mappers, nil
}

// DeleteIgMapper deletes a ig mapper by map_id
func (r *Repository) DeleteIgMapper(mapID string) error {
	_, err := r.db.Model(&IgMapper{
		MapID: uuid.MustParse(mapID),
	}).Where("map_id = ?", mapID).Delete()
	if err != nil {
		return err
	}

	return nil
}

// DeleteIgRunner deletes a ig runner by user_id and map_id
func (r *Repository) DeleteIgRunner(userID, mapID string) error {
	_, err := r.db.Model(&IgRunner{
		UserID: userID,
		MapID:  uuid.MustParse(mapID),
	}).Where("user_id = ? AND map_id = ?", userID, mapID).Delete()
	if err != nil {
		return err
	}

	return nil
}

// SaveIgCaptured saves an IgCaptured to the database
func (r *Repository) SaveIgCaptured(c *IgCaptured) error {
	_, err := r.db.Model(c).Insert()
	if err != nil {	
		return err
	}

	return nil
}

// createTables creates all tables in the database
func (r *Repository) CreateTables() error {
	models := []interface{}{
		(*User)(nil),
		(*Builder)(nil),
		(*EventMapper)(nil),
		(*Event)(nil),
		(*Runner)(nil),
		(*IgMapper)(nil),
		(*IgRunner)(nil),
		(*IgCaptured)(nil),
	}

	for _, model := range models {
		err := r.db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
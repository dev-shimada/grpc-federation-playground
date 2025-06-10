package schema // dev-shimada

import ( // dev-shimada
	"time" // dev-shimada

	"entgo.io/ent"              // dev-shimada
	"entgo.io/ent/schema/field" // dev-shimada
	"github.com/google/uuid"    // dev-shimada
) // dev-shimada

// Message holds the schema definition for the Message entity. // dev-shimada
type Message struct { // dev-shimada
	ent.Schema // dev-shimada
} // dev-shimada

// Fields of the Message. // dev-shimada
func (Message) Fields() []ent.Field { // dev-shimada
	return []ent.Field{ // dev-shimada
		field.UUID("id", uuid.UUID{}). // dev-shimada
						Default(uuid.New). // dev-shimada
						StorageKey("id").  // dev-shimada
						Unique().          // dev-shimada
						Immutable(),       // dev-shimada
		field.String("user_id"). // dev-shimada
						NotEmpty(), // dev-shimada
		field.String("text"). // dev-shimada
					NotEmpty(), // dev-shimada
		field.Time("created_at"). // dev-shimada
						Default(time.Now). // dev-shimada
						Immutable(),       // dev-shimada
		field.Time("updated_at"). // dev-shimada
						Default(time.Now).       // dev-shimada
						UpdateDefault(time.Now), // dev-shimada
	} // dev-shimada
} // dev-shimada

// Edges of the Message. // dev-shimada
func (Message) Edges() []ent.Edge { // dev-shimada
	return nil // dev-shimada
} // dev-shimada

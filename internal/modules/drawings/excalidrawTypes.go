package drawings

type FillStyle string

type StrokeStyle string

type ChartType string

type PointBinding struct {
	ElemntId string  `json:"elementId"`
	Focus    float32 `json:"focus"`
	Gap      float32 `json:"gap"`
}

type PointerType string

type StrokeRoundess string

type GroupId string

type RoundessType float32

type FileId string

type FontFamilyValues float32

type TextAlign string

type VerticalAlign string

type Point [2]float32

type ArrowHead string

// taken from _ExcalidrawElementBase
type ExcalidrawElement struct {
	// MongoId         primitive.ObjectID `json:"_id" bson:"_id"`
	Type            string      `json:"type" bson:"type"`
	Id              string      `json:"id" bson:"id"`
	X               float32     `json:"x" bson:"x"`
	Y               float32     `json:"y" bson:"y"`
	StrokeColor     string      `json:"strokeColor" bson:"strokeColor"`
	BackgroundColor string      `json:"backgroundColor" bson:"backgroundColor"`
	FillStyle       FillStyle   `json:"fillStyle" bson:"fillStyle"`
	StrokeWidth     float32     `json:"strokeWidth" bson:"strokeWidth"`
	StrokeStyle     StrokeStyle `json:"strokeStyle" bson:"strokeStyle"`
	Roundness       struct {
		Type  RoundessType `json:"type" bson:"type"`
		Value float32      `json:"value" bson:"value"`
	} `json:"roundness" bson:"roundness"`
	Roughness float32 `json:"roughness" bson:"roughness"`
	Opacity   float32 `json:"opacity" bson:"opacity"`
	Width     float32 `json:"width" bson:"width"`
	Height    float32 `json:"height" bson:"height"`
	Angle     float32 `json:"angle" bson:"angle"`
	/** Random integer used to seed shape generation so that the roughjs shape
	  doesn't differ across renders. */
	Seed float32 `json:"seed" bson:"seed"`
	/** Integer that is sequentially incremented on each change. Used to reconcile
	  elements during collaboration or when saving to server. */
	Version float32 `json:"version" bson:"version"`
	/** Random integer that is regenerated on each change.
	  Used for deterministic reconciliation of updates during collaboration,
	  in case the versions (see above) are identical. */
	VersionNonce float32 `json:"versionNonce" bson:"versionNonce"`
	IsDeleted    bool    `json:"isDeleted" bson:"isDeleted"`
	// List of groups the element belongs to. Ordered from deepest to shallowest.
	GroupIds []GroupId `json:"groupIds" bson:"groupIds"`
	FrameId  string    `json:"frameId" bson:"frameId"`
	// other elements that are bound to this element
	BoundElements []ExcalidrawElement `json:"boundElements" bson:"boundElements"`
	// epoch (ms) timestamp of last element update
	Updated    float32     `json:"updated" bson:"updated"`
	Link       string      `json:"link" bson:"link"`
	Locked     bool        `json:"locked" bson:"locked"`
	CustomData interface{} `json:"customData" bson:"customData"`

	// taken from ExcaliidrawEmbedableElement
	/**
	 * indicates whether the embeddable src (url) has been validated for rendering.
	 * null value indicates that the validation is pending. We reset the
	 * value on each restore (or url change) so that we can guarantee
	 * the validation came from a trusted source (the editor). Also because we
	 * may not have access to host-app supplied url validator during restore.
	 */
	Validated bool `json:"validated" bson:"validated"`

	// taken from ExcalidrawImageElement
	FileId string `json:"fileId" bson:"fileId"`
	/** Whether respective file is persisted */
	Status string `json:"status" bson:"status"`
	/** X and Y scale factors <-1, 1>, used for image axis flipping */
	Scale [2]float32 `json:"scale" bson:"scale"`

	// taken from ExcalidrawFrameElement
	Name string `json:"name" bson:"name"`
	// taken from ExcalidrawTextElement

	FontSize      float32          `fontSize:"fontSize"`
	FontFamily    FontFamilyValues `json:"fontFamily" bson:"fontFamily"`
	Text          string           `json:"text" bson:"text"`
	Baseline      float32          `json:"baseline" bson:"baseline"`
	TextAlign     TextAlign        `json:"textAlign" bson:"textAlign"`
	VerticalAlign VerticalAlign    `json:"verticalAlign" bson:"verticalAlign"`
	ContainerId   string           `json:"containerId" bson:"containerId"`
	OriginalText  string           `json:"originalText" bson:"originalText"`
	/** Unitless line height (aligned to W3C). To get line height in px, multiply with font size (using `getLineHeightInPx` helper). */
	Linehight float32 `json:"lineHeight" bson:"lineHeight"`

	// taken from ExcalidrawLinearElement
	Points             []Point      `json:"points" bson:"points"`
	LastCommittedPoint Point        `json:"lastCommittedPoint" bson:"lastCommittedPoint"`
	StartBinding       PointBinding `json:"startBinding" bson:"startBinding"`
	EndBinding         PointBinding `json:"endBinding" bson:"endBinding"`
	StartArrowhead     ArrowHead    `json:"startArrowhead" bson:"startArrowhead"`
	EndArrowhead       ArrowHead    `json:"endArrowhead" bson:"endArrowhead"`

	//taken from ExcalidrawFreeDrawElement
	Pressures        []float32 `json:"pressures" bson:"pressures"`
	SimulatePressure bool      `json:"simulatePressure" bson:"simulatePressure"`
}

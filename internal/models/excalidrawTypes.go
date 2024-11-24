package models

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
	Type            string      `json:"type"`
	Id              string      `json:"id"`
	X               float32     `json:"x"`
	Y               float32     `json:"y"`
	StrokeColor     string      `json:"strokeColor"`
	BackgroundColor string      `json:"backgroundColor"`
	FillStyle       FillStyle   `json:"fillStyle"`
	StrokeWidth     float32     `json:"strokeWidth"`
	StrokeStyle     StrokeStyle `json:"strokeStyle"`
	Roundness       struct {
		Type  RoundessType `json:"type"`
		Value float32      `json:"value"`
	} `json:"roundness"`
	Roughness float32 `json:"roughness"`
	Opacity   float32 `json:"opacity"`
	Width     float32 `json:"width"`
	Height    float32 `json:"height"`
	Angle     float32 `json:"angle"`
	/** Random integer used to seed shape generation so that the roughjs shape
	  doesn't differ across renders. */
	Seed float32 `json:"seed"`
	/** Integer that is sequentially incremented on each change. Used to reconcile
	  elements during collaboration or when saving to server. */
	Version float32 `json:"version"`
	/** Random integer that is regenerated on each change.
	  Used for deterministic reconciliation of updates during collaboration,
	  in case the versions (see above) are identical. */
	VersionNonce float32 `json:"versionNonce"`
	IsDeleted    bool    `json:"isDeleted"`
	// List of groups the element belongs to. Ordered from deepest to shallowest.
	GroupIds []GroupId `json:"groupIds"`
	FrameId  string    `json:"frameId"`
	// other elements that are bound to this element
	BoundElements []ExcalidrawElement `json:"boundElements"`
	// epoch (ms) timestamp of last element update
	Updated    float32     `json:"updated"`
	Link       string      `json:"link"`
	Locked     bool        `json:"locked"`
	CustomData interface{} `json:"customData"`

	// taken from ExcaliidrawEmbedableElement
	/**
	 * indicates whether the embeddable src (url) has been validated for rendering.
	 * null value indicates that the validation is pending. We reset the
	 * value on each restore (or url change) so that we can guarantee
	 * the validation came from a trusted source (the editor). Also because we
	 * may not have access to host-app supplied url validator during restore.
	 */
	Validated bool `json:"validated"`

	// taken from ExcalidrawImageElement
	FileId string `json:"fileId"`
	/** Whether respective file is persisted */
	Status string `json:"status"`
	/** X and Y scale factors <-1, 1>, used for image axis flipping */
	Scale [2]float32 `json:"scale"`

	// taken from ExcalidrawFrameElement
	Name string `json:"name"`
	// taken from ExcalidrawTextElement

	FontSize      float32          `fontSize:"fontSize"`
	FontFamily    FontFamilyValues `json:"fontFamily"`
	Text          string           `json:"text"`
	Baseline      float32          `json:"baseline"`
	TextAlign     TextAlign        `json:"textAlign"`
	VerticalAlign VerticalAlign    `json:"verticalAlign"`
	ContainerId   string           `json:"containerId"`
	OriginalText  string           `json:"originalText"`
	/** Unitless line height (aligned to W3C). To get line height in px, multiply with font size (using `getLineHeightInPx` helper). */
	Linehight float32 `json:"lineHeight"`

	// taken from ExcalidrawLinearElement
	Points             []Point      `json:"points"`
	LastCommittedPoint Point        `json:"lastCommittedPoint"`
	StartBinding       PointBinding `json:"startBinding"`
	EndBinding         PointBinding `json:"endBinding"`
	StartArrowhead     ArrowHead    `json:"startArrowhead"`
	EndArrowhead       ArrowHead    `json:"endArrowhead"`

	//taken from ExcalidrawFreeDrawElement
	Pressures        []float32 `json:"pressures"`
	SimulatePressure bool      `json:"simulatePressure"`
}

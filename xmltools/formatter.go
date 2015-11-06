package xmltools

import (
	"bytes"
	"encoding/xml"
	"io"
)

func FormatReader(r io.Reader, prefix, indent string, stripNamespace, stripAttributes bool) (b []byte, err error) {
	buf := new(bytes.Buffer)
	enc := xml.NewEncoder(buf)
	enc.Indent(prefix, indent)

	dec := xml.NewDecoder(r)

	for {
		var t xml.Token
		t, err = dec.Token()
		if t == nil || err != nil {
			if err == io.EOF {
				err = nil
			}

			break
		}

		switch elem := t.(type) {
		case xml.StartElement:
			if stripNamespace {
				elem.Name.Space = ""
			}

			if stripAttributes {
				elem.Attr = nil
			}

			t = elem
		case xml.EndElement:
			if stripNamespace {
				elem.Name.Space = ""
			}

			t = elem
		}

		err = enc.EncodeToken(t)
		if err != nil {
			break
		}
	}
	if err != nil {
		return
	}

	err = enc.Flush()
	if err != nil {
		return
	}

	b = buf.Bytes()
	return
}

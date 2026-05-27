package lnk

import (
	"fmt"
	"io"
)

// File represents one lnk file.
type File struct {
	Header     ShellLinkHeaderSection  // File header.
	IDList     LinkTargetIDListSection // LinkTargetIDList.
	LinkInfo   LinkInfoSection         // LinkInfo.
	StringData StringDataSection       // StringData.
	DataBlocks ExtraDataSection        // ExtraData blocks.
}

// Read parses an io.Reader pointing to the contents of a lnk file.
func Read(r io.Reader, maxSize uint64) (f File, err error) {

	f.Header, err = Header(r, maxSize)
	if err != nil {
		return f, fmt.Errorf("lnk.Read: parse Header - %s", err.Error())
	}

	// If HasLinkTargetIDList is set, header is immediately followed by a LinkTargetIDList.
	if f.Header.LinkFlags["HasLinkTargetIDList"] {
		f.IDList, err = LinkTarget(r)
		if err != nil {
			return f, fmt.Errorf("lnk.Read: parse LinkTarget - %s", err.Error())
		}
	}

	// If HasLinkInfo is set, read LinkInfo section.
	if f.Header.LinkFlags["HasLinkInfo"] {
		f.LinkInfo, err = LinkInfo(r, maxSize)
		if err != nil {
			return f, fmt.Errorf("lnk.Read: parse LinkInfo - %s", err.Error())
		}
	}

	// Read StringData section.
	f.StringData, err = StringData(r, f.Header.LinkFlags)
	if err != nil {
		return f, fmt.Errorf("lnk.Read: parse StringData - %s", err.Error())
	}

	f.DataBlocks, err = DataBlock(r, maxSize)
	if err != nil {
		return f, fmt.Errorf("lnk.Read: parse ExtraDataBlock - %s", err.Error())
	}

	return f, err
}

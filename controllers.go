package main

import "net/http"

func indexController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorReply(w, ErrNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("imaginary server " + Version))
}

func formController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlForm()))
}

func imageControllerDispatcher(o ServerOptions, operation Operation) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf []byte
		var err error

		if r.Method == "GET" && o.Mount != "" {
			buf, err = readLocalImage(w, r, o.Mount)
		} else {
			buf, err = readPayload(w, r)
		}

		if err != nil {
			return
		}

		imageController(w, r, buf, operation)
	}
}

func imageController(w http.ResponseWriter, r *http.Request, buf []byte, Operation Operation) {
	if len(buf) == 0 {
		ErrorReply(w, ErrEmptyPayload)
		return
	}

	mimeType := http.DetectContentType(buf)
	if IsImageMimeTypeSupported(mimeType) == false {
		ErrorReply(w, ErrUnsupportedMedia)
		return
	}

	opts := readParams(r)
	if opts.Type != "" && ImageType(opts.Type) == 0 {
		ErrorReply(w, ErrOutputFormat)
		return
	}

	image, err := Operation.Run(buf, opts)
	if err != nil {
		ErrorReply(w, NewError("Error while processing the image: "+err.Error(), BAD_REQUEST))
		return
	}

	w.Header().Set("Content-Type", image.Mime)
	w.Write(image.Body)
}

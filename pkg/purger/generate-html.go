package purger

import (
	"Stale-purger/pkg/consts"
	"html/template"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (stalePodsInfo StalePodsInfo) GenerateHTMLOutput(logger *logrus.Entry) error {

	// Parse the template
	tmpl, err := template.New("report").Parse(consts.HTMLOutputTemplate)
	if err != nil {
		return errors.Wrap(err, "Couldn't parse the HTML template")
	}

	// Create an HTML file
	file, err := os.Create(consts.HTMLOutputTemplateFileName)
	if err != nil {
		return errors.Wrap(err, "Couldn't create HTML report file")
	}
	defer file.Close()

	// Execute the template with data and write to file
	err = tmpl.Execute(file, stalePodsInfo)
	if err != nil {
		return errors.Wrap(err, "Couldn't write stale pod information into HTML template")
	}

	logger.Infof("HTML file successfully generated: %s", consts.HTMLOutputTemplateFileName)
	return nil
}

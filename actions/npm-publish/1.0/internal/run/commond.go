package run

import "github.com/erda-project/erda-actions/actions/npm-publish/1.0/internal/npm"

type Command struct {
	packageManager npm.PackageManager
}

func NewCommand(packageManager npm.PackageManager) *Command {
	return &Command{
		packageManager: packageManager,
	}
}

func (command *Command) Run(request Request) (Response, error) {
	//Install npm-cli-login to use login command.
	err := command.packageManager.Install(
		"npm-cli-login",
		"",
		true,
	)
	if err != nil {
		return Response{}, err
	}

	err = command.packageManager.Login(
		request.Params.UserName,
		request.Params.Password,
		request.Params.Email,
		request.Source.Registry,
	)
	if err != nil {
		return Response{}, err
	}

	err = command.packageManager.Publish(
		request.Params.Path,
		request.Params.Tag,
		request.Source.Registry,
	)
	if err != nil {
		return Response{}, err
	}

	out, err := command.packageManager.View(
		request.Source.PackageName,
		request.Source.Registry,
	)
	if err != nil {
		return Response{}, err
	}

	err = command.packageManager.Logout(
		request.Source.Registry,
	)
	if err != nil {
		return Response{}, err
	}

	return Response{
		Version: resource.Version{
			Version: out.Version,
		},
		Metadata: []resource.MetadataPair{
			{
				Name:  "name",
				Value: out.Name,
			},
		},
	}, nil
}

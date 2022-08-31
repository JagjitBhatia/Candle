package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	_ "database/sql"
	"fmt"
	"strconv"

	"github.com/JagjitBhatia/Candle/graph/generated"
	"github.com/JagjitBhatia/Candle/graph/model"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func (r *mutationResolver) CreateUser(ctx context.Context, username string, firstName string, lastName string, institution string, pfpURL *string) (*model.User, error) {

	var newUser model.User

	newUser.Username = username
	newUser.FirstName = firstName
	newUser.LastName = lastName
	newUser.Institution = institution

	if pfpURL != nil {
		newUser.PfpURL = pfpURL
	} else {
		*(newUser.PfpURL) = "<default_profile_pic_url>" // TODO: Replace with actual default pic url
	}

	createUserQuery := fmt.Sprintf("INSERT INTO Users VALUES(NULL,'%s', '%s', '%s', '%s', '%s')",
		newUser.Username,
		newUser.FirstName,
		newUser.LastName,
		newUser.Institution,
		*newUser.PfpURL,
	)

	result, err := r.Db.Exec(createUserQuery)

	if err != nil {
		log.Errorf("failed to create user with error %v", err)
		return nil, err
	}

	newUserID, err := result.LastInsertId()

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	newUser.ID = strconv.FormatInt(newUserID, 10)

	return &newUser, nil
}

func (r *mutationResolver) CreateOrg(ctx context.Context, name string, institution string, orgPicURL *string, userID string, title string) (*model.Org, error) {
	var newOrg model.Org
	var newMember model.Member
	var user model.User

	user.ID = userID

	findUserQuery := fmt.Sprintf("SELECT username, first_name, last_name, institution, pfp_url FROM Users WHERE id = %s", user.ID)

	results, err := r.Db.Query(findUserQuery)

	if err != nil {
		log.Errorf("failed to create org with error: %v", err)
		return nil, err
	}

	userExists := false

	for results.Next() {
		err = results.Scan(&user.Username, &user.FirstName, &user.LastName, &user.Institution, &user.PfpURL)
		if err != nil {
			log.Error(err)
			continue
		}

		userExists = true
	}

	if !userExists {
		err = fmt.Errorf("failed to create org with error: user not found")
		log.Error(err)
		return nil, err
	}

	newOrg.Name = name
	newOrg.Institution = institution

	if orgPicURL != nil {
		newOrg.OrgPicURL = orgPicURL
	} else {
		*(newOrg.OrgPicURL) = "<default_profile_pic_url>" // TODO: Replace with actual default pic url
	}

	createOrgQuery := fmt.Sprintf("INSERT INTO Orgs VALUES(NULL, '%s', '%s', '%s')", newOrg.Name, newOrg.Institution, *newOrg.OrgPicURL)

	result, err := r.Db.Exec(createOrgQuery)

	if err != nil {
		log.Errorf("failed to create org with error: %v", err)
		return nil, err
	}

	newOrgId, err := result.LastInsertId()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	newOrg.ID = strconv.FormatInt(newOrgId, 10)

	newMember.User = &user
	newMember.Role = "admin"
	newMember.Title = title

	addOrgMemberQuery := fmt.Sprintf("INSERT INTO Members VALUES (%s, %s, '%s', '%s')", user.ID, newOrg.ID, newMember.Role, newMember.Title)

	_, err = r.Db.Exec(addOrgMemberQuery)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	newOrg.Members = append(newOrg.Members, &newMember)

	return &newOrg, err
}

func (r *mutationResolver) AddOrgMember(ctx context.Context, newMemberID string, orgID string, role string, title string) (*model.Org, error) {
	var org model.Org
	var newMember model.Member
	var user model.User

	org.ID = orgID
	user.ID = newMemberID

	newMember.Role = role
	newMember.Title = title

	findOrgQuery := fmt.Sprintf("SELECT org_name, institution, org_pic_url FROM Orgs WHERE id = %s", org.ID)
	findOrgMembers := fmt.Sprintf(`SELECT Users.id, Users.username, Users.first_name, Users.last_name, Users.institution, Users.pfp_url, Members.member_role, Members.title
									FROM Users JOIN Members ON Users.id = Members.user_id WHERE Members.org_id = %s`, org.ID)
	findUserQuery := fmt.Sprintf("SELECT username, first_name, last_name, institution, pfp_url FROM Users WHERE id = %s", user.ID)
	addOrgMemberQuery := fmt.Sprintf("INSERT INTO Members VALUES (%s, %s, '%s', '%s')", user.ID, org.ID, newMember.Role, newMember.Title)

	results, err := r.Db.Query(findOrgQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	orgExists := false

	for results.Next() {
		err = results.Scan(&org.Name, &org.Institution, &org.OrgPicURL)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		orgExists = true
	}

	if !orgExists {
		err = fmt.Errorf("error: org not found")
		log.Error(err)
		return nil, err
	}

	results, err = r.Db.Query(findOrgMembers)

	if err != nil {
		log.Error(err.Error())
	}

	for results.Next() {
		var currentMember model.Member
		var memberUser model.User

		err = results.Scan(&memberUser.ID, &memberUser.Username, &memberUser.FirstName, &memberUser.LastName, &memberUser.Institution, &memberUser.PfpURL,
			&currentMember.Role, &currentMember.Title)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		currentMember.User = &memberUser
		org.Members = append(org.Members, &currentMember)
	}

	results, err = r.Db.Query(findUserQuery)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	userExists := false

	for results.Next() {
		err = results.Scan(&user.Username, &user.FirstName, &user.LastName, &user.Institution, &user.PfpURL)

		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		userExists = true
	}

	if !userExists {
		err = fmt.Errorf("error: user not found")
		return nil, err
	}

	newMember.User = &user

	_, err = r.Db.Exec(addOrgMemberQuery)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	org.Members = append(org.Members, &newMember)

	return &org, err
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	var users []*model.User

	results, err := r.Db.Query("SELECT id, username, first_name, last_name, institution, pfp_url FROM Users")
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	for results.Next() {
		var user model.User

		err = results.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Institution, &user.PfpURL)

		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *queryResolver) Orgs(ctx context.Context) ([]*model.Org, error) {
	var orgs []*model.Org

	results, err := r.Db.Query("SELECT id, org_name, institution, org_pic_url FROM Orgs")
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	for results.Next() {
		var org model.Org

		err = results.Scan(&org.ID, &org.Name, &org.Institution, &org.OrgPicURL)

		if err != nil {
			log.Error(err.Error())
			continue
		}

		results, err := r.Db.Query(fmt.Sprintf(`SELECT Users.id, Users.username, Users.first_name, Users.last_name, Users.institution, Users.pfp_url,
												Members.member_role, Members.title FROM Users JOIN Members ON Users.id = Members.user_id
												WHERE Members.org_id = %s`, org.ID))
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		for results.Next() {
			var member model.Member
			var user model.User

			err = results.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Institution, &user.PfpURL,
				&member.Role, &member.Title)

			if err != nil {
				log.Error(err.Error())
				continue
			}

			member.User = &user
			org.Members = append(org.Members, &member)
		}

		orgs = append(orgs, &org)
	}

	return orgs, nil
}

func (r *queryResolver) UserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	user.ID = id

	findUserQuery := fmt.Sprintf("SELECT username, first_name, last_name, institution, pfp_url FROM Users WHERE id = %s", user.ID)
	results, err := r.Db.Query(findUserQuery)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	for results.Next() {
		err = results.Scan(&user.Username, &user.FirstName, &user.LastName, &user.Institution, &user.PfpURL)

		if err != nil {
			panic(err.Error())
		}
	}

	return &user, nil
}

func (r *queryResolver) UserByName(ctx context.Context, name string) ([]*model.User, error) {
	var users []*model.User

	results, err := r.Db.Query("SELECT id, username, first_name, last_name, institution, pfp_url FROM Users WHERE username = '%s'", name)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	for results.Next() {
		var user model.User

		err = results.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Institution, &user.PfpURL)

		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *queryResolver) OrgByID(ctx context.Context, id string) (*model.Org, error) {
	var org model.Org
	org.ID = id

	findOrgQuery := fmt.Sprintf("SELECT org_name, institution, org_pic_url FROM Orgs WHERE id = %s", org.ID)
	results, err := r.Db.Query(findOrgQuery)

	if err != nil {
		err = fmt.Errorf("error running orgById query. Error: %v", err)
		log.Error(err)
		return nil, err

	}

	for results.Next() {
		err = results.Scan(&org.Name, &org.Institution, &org.OrgPicURL)

		if err != nil {
			log.Error(err.Error())
			continue
		}

	}

	results, err = r.Db.Query(fmt.Sprintf(`SELECT Users.id, Users.username, Users.first_name, Users.last_name, Users.institution, Users.pfp_url,
												Members.member_role, Members.title FROM Users JOIN Members ON Users.id = Members.user_id
												WHERE Members.org_id = %s`, org.ID))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	for results.Next() {
		var member model.Member
		var user model.User

		err = results.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Institution, &user.PfpURL,
			&member.Role, &member.Title)

		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		member.User = &user
		org.Members = append(org.Members, &member)
	}

	return &org, nil
}

func (r *queryResolver) OrgByName(ctx context.Context, name string) ([]*model.Org, error) {
	var orgs []*model.Org

	results, err := r.Db.Query(fmt.Sprintf("SELECT id, org_name, institution, org_pic_url FROM Orgs WHERE org_name = '%s'", name))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	for results.Next() {
		var org model.Org

		err = results.Scan(&org.ID, &org.Name, &org.Institution, &org.OrgPicURL)

		if err != nil {
			panic(err.Error())
		}

		results, err := r.Db.Query(fmt.Sprintf(`SELECT Users.id, Users.username, Users.first_name, Users.last_name, Users.institution, Users.pfp_url,
												Members.member_role, Members.title FROM Users JOIN Members ON Users.id = Members.user_id
												WHERE Members.org_id = %s`, org.ID))
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		for results.Next() {
			var member model.Member
			var user model.User

			err = results.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Institution, &user.PfpURL,
				&member.Role, &member.Title)

			if err != nil {
				log.Error(err.Error())
				return nil, err
			}

			member.User = &user
			org.Members = append(org.Members, &member)
		}

		orgs = append(orgs, &org)
	}

	return orgs, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

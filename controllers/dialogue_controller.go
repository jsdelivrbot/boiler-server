package controllers

import (
	"github.com/AzureRelease/boiler-server/dba"
	"fmt"
	"github.com/AzureRelease/boiler-server/models"
	"encoding/json"
	"github.com/AzureTech/goazure/orm"
	"github.com/AzureTech/goazure"
)

type DialogueController struct {
	MainController
}

func (ctl *DialogueController) DialogueList() {
	usr := ctl.GetCurrentUser()

	var dialogues []models.Dialogue
	var aDialogues []models.Dialogue
	qs := dba.BoilerOrm.QueryTable("dialogue")
	qs = qs.RelatedSel("CreatedBy__Organization").RelatedSel("CreatedBy__Role")
	if usr.IsCommonUser() ||
		usr.Status == models.USER_STATUS_INACTIVE || usr.Status == models.USER_STATUS_NEW {
		qs = qs.Filter("IsDemo", true)
	} else {
		qs = qs.Filter("IsDemo", false)
		if usr.IsOrganizationUser() {
			qs = qs.Filter("CreatedBy__Organization__Uid", usr.Organization.Uid)
		}
	}
	num, err := qs.Filter("IsDeleted", false).OrderBy("-UpdatedDate").All(&dialogues)
	fmt.Printf("Returned Rows Num: %d, %s", num, err)
	for i, d := range dialogues {
		var comments []*models.DialogueComment
		qc := dba.BoilerOrm.QueryTable("dialogue_comment")
		qc = qc.RelatedSel("Dialogue").RelatedSel("CreatedBy__Organization").RelatedSel("CreatedBy__Role")
		qc = qc.Filter("Dialogue__Uid", d.Uid)
		num, err := qc.Filter("IsDeleted", false).OrderBy("CreatedDate").All(&comments)
		d.Comments = comments
		aDialogues = append(aDialogues, d)
		fmt.Println("[",i,"]", d, num, err)
	}


	ctl.Data["json"] = aDialogues
	ctl.ServeJSON()
}

type aComment struct {
	Uid			string		`json:"uid"`
	DialogueId 		string		`json:"dialogueId"`
	Topic			string		`json:"topic"`
	Content			string		`json:"content"`
	Attachment		string		`json:"attachment"`
}

func (ctl *DialogueController) CommentUpdate() {
	usr := ctl.GetCurrentUser()
	var cmt aComment
	comment := models.DialogueComment{}

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &cmt); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal Comment Error", err)
		return
	}

	if len(cmt.Uid) > 0 {
		fmt.Println("\n>>>>>>There is Comment Uid:", cmt.Uid)
		comment.Uid = cmt.Uid

		if err := DataCtl.ReadData(&comment); err != nil {
			fmt.Println("Read Exist Comment Error:", err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("Read Exist AlarmRule Parameter Error:"))
			return
		}
	}

	dialogue := models.Dialogue{}

	qs := dba.BoilerOrm.QueryTable("dialogue")
	qs = qs.RelatedSel("CreatedBy__Organization")
	qs = qs.Filter("Uid", cmt.DialogueId)
	if err := qs.Filter("IsDeleted", false).One(&dialogue); err != nil {
		fmt.Println("Read Dialogue Error:", err)
		if err == orm.ErrNoRows || err == orm.ErrMissPK {
			d := models.Dialogue{}
			var id int64 = 0
			if err := dba.BoilerOrm.QueryTable("dialogue").OrderBy("-DialogueId").One(&d, "DialogueId"); err != nil {
				fmt.Println("Read Dialogue LastId Error:", err)
			} else {
				id = d.DialogueId + 1;
			}
			dialogue.DialogueId = id;
			dialogue.Name = cmt.Topic
			dialogue.CreatedBy = usr;
			dialogue.UpdatedBy = usr;
			dialogue.Status = 1;
			if err := DataCtl.AddData(&dialogue, true); err != nil {
				fmt.Println("Add Dialogue Error:", err)
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("Update Comment Error!"))
				return
			}
		}
	} else {
		dialogue.UpdatedBy = usr;
		if usr.Organization != nil && dialogue.CreatedBy.Organization != nil &&
			usr.Organization.Uid == dialogue.CreatedBy.Organization.Uid {
			dialogue.Status = 2;
		} else {
			dialogue.Status = 3;
		}

	}

	comment.Dialogue = &dialogue

	comment.Name = cmt.Topic
	comment.Content = cmt.Content
	comment.Attachment = cmt.Attachment
	comment.From = usr;
	comment.CreatedBy = usr;
	comment.UpdatedBy = usr;

	if err := DataCtl.AddData(&comment, true); err != nil {
		fmt.Println("Add Comment Error:", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Update Comment Error!"))
		return
	}

	if err := DataCtl.UpdateData(&dialogue); err != nil {
		fmt.Println("Update Dialogue Error:", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Update Dialogue Error!"))
		return
	}

	fmt.Println("\nUpdate Comment:", comment, cmt)
}

func (ctl *DialogueController) DialogueDelete() {
	usr := ctl.GetCurrentUser()
	var cmt aComment
	var dialogue models.Dialogue

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &cmt); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Comment Error", err)
		return
	}

	if err := dba.BoilerOrm.QueryTable("dialogue").RelatedSel("CreatedBy__Role").Filter("Uid", cmt.Uid).One(&dialogue); err != nil {
		e := fmt.Sprintf("Read Dialogue for Delete Error: %v", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	author := dialogue.CreatedBy

	goazure.Error("Usr:", usr, "\n===========\n", "Author:", author)

	if 	author == nil ||
		(usr.IsAdmin() && usr.Role.RoleId < author.Role.RoleId) ||
		author.Uid == usr.Uid ||
		(usr.IsOrganizationUser() && usr.Organization.Uid == author.Organization.Uid && usr.Role.RoleId == models.USER_ROLE_ORG_ADMIN) {
		if err := DataCtl.DeleteData(&dialogue); err != nil {
			e := fmt.Sprintf("Delete Dialogue Error: %v", err)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}
	}
}
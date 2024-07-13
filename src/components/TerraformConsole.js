import React, {useState, useEffect} from 'react';
import axios from 'axios';

export default function TerraformConsole() {
    const url = '<ip/url of backend application>'
    const [applyButton, setApplyButton] = useState(true);
    const [output, setOutput] = useState();
    const [text, setText] = useState();
    const [alert, setAlert] = useState({message: null, type: null});

    useEffect(() => {
        const savedText = localStorage.getItem('terraformCode');
        if (savedText) {
            setText(savedText);
        } else {
            axios.get(url+'/save')
                .then((response) => {
                    console.log(response);
                    setText(response.data);
                })
                .catch((error) => {
                    console.log(error);
                });
        }
    }, []);

    const showAlert = (message, type) => {
        setAlert({message, type});
        //if (type === "success") {
        //    setTimeout(() => {
        //        setAlert({message: null, type: null});
        //    }, 2000);
        //}
    }

    const onClickSave = (e) => {
	showAlert("Saving...", "warning");
	e.target.disabled = true;
        const data = text;
        axios.post(url+'/save', {
            data}).then((response) => {
                console.log(response);
                localStorage.removeItem('terraformCode');
                // alert('Terraform code saved successfully!');
                showAlert("Saved successfully", "success");
		e.target.disabled = false;
		setApplyButton(true);
            }).catch((error) => {
                console.log(error);
		showAlert("Error in saving file", "danger");
		e.target.disabled = false;
            });
    }
    const onClickInit = (e) => {
        //showAlert("Terraform init in progress", "warning");
        e.target.disabled = true;
	showAlert("Terraform init in progress", "warning");
        axios.post(url+'/init', {
        }).then((response) => {
            // console.log(response);
            setOutput(response.data);
            //alert('Terraform initialized successfully');
            showAlert("Terraform init successful", "success");
            e.target.disabled = false;
        }).catch((error) => {
            console.log(error);
            showAlert("Terraform init failed", "danger");
            e.target.disabled = false;
        });
    }
    const onClickPlan = (e) => {
        //showAlert("Terraform plan in progress", "warning");
        e.target.disabled = true;
	showAlert("Terraform plan in progress", "warning");
        axios.post(url+'/plan', {
        }).then((response) => {
            // console.log(response);
            setOutput(response.data);
            //alert('Terraform plan successful');
            showAlert("Terraform plan successful", "success");
            e.target.disabled = false;
	    setApplyButton(false);
            // document.getElementsByName('apply')[0].disabled = false;
        }).catch((error) => {
            console.log(error);
            showAlert("Terraform plan failed", "danger");
            e.target.disabled = false;
        });
    }
    const onClickApply = (e) => {
        //showAlert("Terraform apply in progress", "warning");
        e.target.disabled = true;
	showAlert("Terraform apply in progress", "warning");
        axios.post(url+'/apply', {
        }).then((response) => {
            // console.log(response);
            setOutput(response.data);
            //alert('Terraform apply successful');
            showAlert("Terraform plan successful", "success");
	    setApplyButton(true);
            e.target.disabled = false;
        }).catch((error) => {
            console.log(error);
            showAlert("Terraform plan failed", "danger");
            e.target.disabled = false;
        });
    }
    const handleOnChange = (e) => {
        setText(e.target.value);
        localStorage.setItem('terraformCode', e.target.value);
        console.log("on change");
    }
        
    return (
        <>
        <div className={`alert alert-${alert.type} alert-dismissible fade show`} role="alert">
          {alert.message}
        </div>
        <div className="mb-3">
            <label htmlFor="terraformCode" className="form-label"><h5>My Application 1</h5></label>
            <textarea className="form-control" id="terraformCode" rows="10" name="terraformCode" onChange={handleOnChange} placeholder="Add terraform code here" value={text}></textarea><br/>
            <button type="button" className="btn btn-dark col-lg-2" onClick={onClickSave} name="save">Save</button><br/><br/>
            <button type="button" className="btn btn-dark col-lg-3" onClick={onClickInit} name="init">Terraform init</button><br/><br/>
            <button type="button" className="btn btn-dark col-lg-3" onClick={onClickPlan} name="plan">Terraform plan</button><br/><br/>
            <button type="button" className="btn btn-dark col-lg-3" onClick={onClickApply} name="apply" disabled={applyButton}>Terraform apply</button><br/><br/>
        </div>
        <div className="mb-3">
            <label htmlFor="terraformOutput" className="form-label"><h6>Output</h6></label>
            <textarea className="form-control" id="terraformOutput" rows="10" name="terraformOutput" placeholder="Terraform output here" value={output} readOnly></textarea>
        </div>
        </>
    )
}


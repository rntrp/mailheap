"use strict";
const limit = 10;
const files = [];
var lastId = 0;
var currentId = 0;
var currentEml = null;
function previewMail(id) {
    const previewHtml = document.getElementById("preview-html");
    previewHtml.classList.add("hidden");
    previewHtml.removeAttribute("title");
    previewHtml.removeAttribute("srcdoc");
    const previewPlain = document.getElementById("preview-plain");
    previewPlain.classList.add("hidden");
    previewPlain.textContent = null;
    const previewHeaders = document.getElementById("preview-headers");
    previewHeaders.classList.add("hidden");
    previewHeaders.textContent = null;
    const title = document.getElementById("preview-title");
    title.textContent = null;
    const subtitle = document.getElementById("preview-subtitle");
    subtitle.textContent = null
    fetch("/mail/" + id)
        .then(response => response.text())
        .then(eml => {
            currentId = id;
            currentEml = eml;
            const parsed = jsmimeparser.parseMail(eml);
            const inlined = inlineImages(parsed.body.html, parsed.attachments);
            previewHtml.title = parsed.subject;
            previewHtml.srcdoc = inlined;
            previewPlain.textContent = parsed.body.text;
            const headerDelim = eml.search(/(\r?\n){2}/g);
            previewHeaders.textContent = eml.substring(0, headerDelim);
            title.textContent = parsed.subject;
            subtitle.textContent = "From " +
                [parsed.from?.email, parsed.date?.toISOString()].join(" at ");
            fileAttachments(parsed.attachments);
            if (parsed.body.html) {
                previewHtml.classList.remove("hidden");
            } else if (parsed.body.text) {
                previewPlain.classList.remove("hidden");
            } else {
                previewHeaders.classList.remove("hidden");
            }
        });
}
function showHtml() {
    document.getElementById("preview-html").classList.remove("hidden");
    document.getElementById("preview-plain").classList.add("hidden");
    document.getElementById("preview-headers").classList.add("hidden");
}
function showPlain() {
    document.getElementById("preview-html").classList.add("hidden");
    document.getElementById("preview-plain").classList.remove("hidden");
    document.getElementById("preview-headers").classList.add("hidden");
}
function showHeaders() {
    document.getElementById("preview-html").classList.add("hidden");
    document.getElementById("preview-plain").classList.add("hidden");
    document.getElementById("preview-headers").classList.remove("hidden");
}
function downloadEml() {
    if (currentEml) {
        const a = document.createElement("a");
        a.download = currentId + ".eml";
        a.href = URL.createObjectURL(new Blob([currentEml], {type: "message/rfc822"}));
        a.rel = "noreferrer noopener nofollow";
        a.target = "_blank"
        try {
            a.click();
        } finally {
            URL.revokeObjectURL(a.href);
        }
    }
}
async function loadMails() {
    const response = await fetch("/mails/" + lastId + "?limit=" + limit);
    const result = await response.json();
    for (const mail of result.data) {
        lastId = mail.id;
        addEmailToList(mail.id, mail.from, mail.to, mail.subject);
    }
    document.getElementById("mail-count").textContent = `(${result.total})`;
    return result;
}
function infiniteScroll() {
    const list = document.getElementById("list");
    const last = list.lastElementChild?.clientHeight ?? 0;
    if (list.scrollTop + list.clientHeight + last >= list.scrollHeight) {
        loadMails().then(result => {
            if (!result.data || result.data.length < limit) {
                list.onscrollend = null;
            }
        });
    }
}
function addEmailToList(id, from, to, subject) {
    const outer = document.createElement("div");
    outer.id = id;
    outer.tabIndex = 8192;
    outer.className = "mail-item";
    outer.onfocus = () => previewMail(id);
    const emailFromText = JSON.parse(from).join("; ");
    const emailFrom = document.createElement("h5");
    emailFrom.className = "mail-from";
    emailFrom.textContent = emailFromText;
    emailFrom.title = emailFromText;
    const emailToText = JSON.parse(to).join("; ");
    const emailTo = document.createElement("h5");
    emailTo.className = "mail-to";
    emailTo.textContent = emailToText;
    emailTo.title = emailToText;
    const emailSubject = document.createElement("h4");
    emailSubject.className = "mail-subject"
    emailSubject.textContent = subject;
    const info = document.createElement("div");
    info.appendChild(emailFrom);
    info.appendChild(emailTo);
    info.appendChild(emailSubject);
    outer.appendChild(info);
    document.getElementById("list").appendChild(outer);
}
function fileAttachments(attachments) {
    while (files.length > 0) {
        URL.revokeObjectURL(files.pop().blob);
    }
    for (const att of attachments) {
        if (att && att.content && att.contentDisposition !== "inline") {
            files.push({
                blob: URL.createObjectURL(new Blob([att.content], { type: att.contentType })),
                name: att.fileName,
                size: att.content.length,
                type: att.contentType,
            });
        }
    }
    const footer = document.getElementById("footer-attachments");
    footer.replaceChildren();
    footer.classList.add("hidden");
    if (files.length > 0) {
        const formatBytes = bytes => {
            if (bytes >= 1048576) {
                return (bytes / 1048576).toFixed(2) + "\xa0MiB";
            } else if (bytes >= 1024) {
                return (bytes / 1024).toFixed(2) + "\xa0KiB";
            } else {
                return bytes + "\xa0B";
            }
        }
        for (const file of files) {
            const a = document.createElement("a");
            a.download = file.name;
            a.href = file.blob;
            a.rel = "noreferrer noopener nofollow";
            a.target = "_blank";
            a.textContent = file.name;
            const size = document.createElement("span");
            size.classList.add("mail-attachment-size");
            size.textContent = "\xa0(" + formatBytes(file.size) + ")";
            const div = document.createElement("div");
            div.classList.add("mail-attachment")
            div.appendChild(a);
            div.appendChild(size);
            footer.appendChild(div);
        }
        footer.classList.remove("hidden");
    }
}
function inlineImages(htmlString, attachments) {
    const inlineMap = new Map();
    for (const att of attachments) {
        if (att.contentId && att.contentDisposition === "inline"
            && ["image/gif", "image/jpeg", "image/png"].includes(att.contentType)) {
            inlineMap.set(att.contentId, att);
        }
    }
    const html = document.createElement("html");
    html.innerHTML = htmlString;
    for (const img of html.getElementsByTagName("img")) {
        if (img.src && img.src.length > 4 && img.src.substring(0, 4).toLowerCase() === "cid:") {
            const cid = "<" + img.src.substring(4) + ">";
            const att = inlineMap.get(cid);
            if (att && att.content) {
                const b64 = btoa(String.fromCharCode.apply(null, att.content));
                img.setAttribute("src", "data:" + att.contentType + ";base64," + b64);
            }
        }
    }
    return html.outerHTML;
}
function uploadMail(eml) {
    const formData = new FormData();
    formData.append("eml", eml);
    fetch("/upload", { method: "POST", body: formData })
        .then(_ => window.location.reload(), err => alert("Upload failed: " + err));
}
function deleteAllMails() {
    if (confirm("Delete all mails?")) {
        fetch("/mails", { method: "DELETE" })
            .then(_ => {
                window.location.reload();
                const previewHtml = document.getElementById("preview-html");
                previewHtml.classList.remove("hidden");
                previewHtml.removeAttribute("title");
                previewHtml.removeAttribute("srcdoc");
            }, err => alert("Delete all mails failed: " + err));
    }
}
window.onload = function() {
    document.getElementById("list").scrollTop = 0;
    loadMails().then(_ => document.querySelector("#list > div:first-child")?.focus());
}
document.getElementById("inbox").onclick = () => window.location.reload();
document.getElementById("upload").onchange = (e) => uploadMail(e.target.files[0]);
document.getElementById("uploadLink").onclick = () => document.getElementById("upload").click();
document.getElementById("delete").onclick = deleteAllMails;
document.getElementById("showHtml").onclick = showHtml;
document.getElementById("showPlain").onclick = showPlain;
document.getElementById("showHeaders").onclick = showHeaders;
document.getElementById("downloadEml").onclick = downloadEml;
document.getElementById("list").onscrollend = infiniteScroll;

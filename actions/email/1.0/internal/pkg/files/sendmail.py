import smtplib
import os
from email.mime.text import MIMEText
from email.header import Header

smtpObj = smtplib.SMTP(os.getenv("python_smtp_host"))
smtpObj.sendmail(os.getenv("python_smtp_sender"), os.getenv("python_smtp_receivers").split(","), os.getenv("python_smtp_message"))

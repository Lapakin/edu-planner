#!/bin/bash
awslocal sqs create-queue --queue-name user-management-queue
awslocal sqs create-queue --queue-name syllabus-message-queue

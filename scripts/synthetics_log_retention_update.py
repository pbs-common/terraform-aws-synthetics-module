import boto3
from botocore.exceptions import ClientError


def get_matching_log_groups(logs_client):
    paginator = logs_client.get_paginator("describe_log_groups")
    matching_logs = list()

    for page in paginator.paginate():
        for log_group in page.get("logGroups", []):
            log_group_name = log_group.get("logGroupName", "")

            if "cwsyn" in log_group_name:
                matching_logs.append(log_group)

    return matching_logs


def lambda_handler(event, context):
    logs_client = boto3.client("logs")

    for log_group in get_matching_log_groups(logs_client)[:1]:
        log_group_name = log_group.get("logGroupName")

        try:
            logs_client.put_retention_policy(
                logGroupName=log_group_name,
                retentionInDays=30,
            )
        except ClientError:
            continue

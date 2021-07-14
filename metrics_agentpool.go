package main

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	devopsClient "github.com/webdevops/azure-devops-exporter/azure-devops-client"
	prometheusCommon "github.com/webdevops/go-prometheus-common"
)

type MetricsCollectorAgentPool struct {
	CollectorProcessorAgentPool

	prometheus struct {
		agentPool               *prometheus.GaugeVec
		agentPoolSize           *prometheus.GaugeVec
		agentPoolUsage          *prometheus.GaugeVec
		agentPoolAgent          *prometheus.GaugeVec
		agentPoolAgentStatus    *prometheus.GaugeVec
		agentPoolAgentJob       *prometheus.GaugeVec
		agentPoolQueueLength    *prometheus.GaugeVec
		agentPoolJobRequestInfo *prometheus.GaugeVec
		agentPoolJobAssignTime  *prometheus.GaugeVec
		agentPoolJobQueueTime   *prometheus.GaugeVec
		agentPoolJobReceiveTime *prometheus.GaugeVec
		agentPoolJobFinishTime  *prometheus.GaugeVec
	}
}

func (m *MetricsCollectorAgentPool) Setup(collector *CollectorAgentPool) {
	m.CollectorReference = collector

	m.prometheus.agentPool = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_devops_agentpool_info",
			Help: "Azure DevOps agentpool",
		},
		[]string{
			"agentPoolID",
			"agentPoolName",
			"agentPoolType",
			"isHosted",
		},
	)
	prometheus.MustRegister(m.prometheus.agentPool)

	m.prometheus.agentPoolSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_devops_agentpool_size",
			Help: "Azure DevOps agentpool",
		},
		[]string{
			"agentPoolID",
		},
	)
	prometheus.MustRegister(m.prometheus.agentPoolSize)

	m.prometheus.agentPoolUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_devops_agentpool_usage",
			Help: "Azure DevOps agentpool usage",
		},
		[]string{
			"agentPoolID",
		},
	)
	prometheus.MustRegister(m.prometheus.agentPoolUsage)

	m.prometheus.agentPoolAgent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_devops_agentpool_agent_info",
			Help: "Azure DevOps agentpool",
		},
		[]string{
			"agentPoolID",
			"agentPoolAgentID",
			"agentPoolAgentName",
			"agentPoolAgentVersion",
			"provisioningState",
			"maxParallelism",
			"agentPoolAgentOs",
			"enabled",
			"status",
			"hasAssignedRequest",
		},
	)
	prometheus.MustRegister(m.prometheus.agentPoolAgent)

	m.prometheus.agentPoolAgentStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_devops_agentpool_agent_status",
			Help: "Azure DevOps agentpool",
		},
		[]string{
			"agentPoolAgentID",
			"type",
		},
	)
	prometheus.MustRegister(m.prometheus.agentPoolAgentStatus)

	m.prometheus.agentPoolAgentJob = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_devops_agentpool_agent_job",
			Help: "Azure DevOps agentpool",
		},
		[]string{
			"agentPoolAgentID",
			"jobRequestId",
			"definitionID",
			"definitionName",
			"planType",
			"scopeID",
		},
	)
	prometheus.MustRegister(m.prometheus.agentPoolAgentJob)

	m.prometheus.agentPoolQueueLength = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_devops_agentpool_queue_length",
			Help: "Azure DevOps agentpool",
		},
		[]string{
			"agentPoolID",
		},
	)
	prometheus.MustRegister(m.prometheus.agentPoolQueueLength)

	m.prometheus.agentPoolJobRequestInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_devops_agentpool_job_request_info",
			Help: "Azure Devops Agentpool",
		},
		[]string{
			"Name",
			"jobRequestID",
			"JobID",
			"Result",
		},
	)
	prometheus.MustRegister(m.prometheus.agentPoolJobRequestInfo)

	m.prometheus.agentPoolJobAssignTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_devops_agentpool_job_assigntime",
			Help: "Azure Devops Agentpool",
		},
		[]string{
			"jobRequestID",
		},
	)
	prometheus.MustRegister(m.prometheus.agentPoolJobAssignTime)

	m.prometheus.agentPoolJobQueueTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_devops_agentpool_job_queuetime",
			Help: "Azure Devops Agentpool",
		},
		[]string{
			"jobRequestID",
		},
	)
	prometheus.MustRegister(m.prometheus.agentPoolJobQueueTime)

	m.prometheus.agentPoolJobReceiveTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_devops_agentpool_job_receivetime",
			Help: "Azure Devops Agentpool",
		},
		[]string{
			"jobRequestID",
		},
	)
	prometheus.MustRegister(m.prometheus.agentPoolJobReceiveTime)

	m.prometheus.agentPoolJobFinishTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_devops_agentpool_job_finishtime",
			Help: "Azure Devops Agentpool",
		},
		[]string{
			"jobRequestID",
		},
	)
	prometheus.MustRegister(m.prometheus.agentPoolJobFinishTime)
}

func (m *MetricsCollectorAgentPool) Reset() {
	m.prometheus.agentPool.Reset()
	m.prometheus.agentPoolSize.Reset()
	m.prometheus.agentPoolAgent.Reset()
	m.prometheus.agentPoolAgentStatus.Reset()
	m.prometheus.agentPoolAgentJob.Reset()
	m.prometheus.agentPoolQueueLength.Reset()
	m.prometheus.agentPoolJobRequestInfo.Reset()
	m.prometheus.agentPoolJobAssignTime.Reset()
	m.prometheus.agentPoolJobQueueTime.Reset()
	m.prometheus.agentPoolJobReceiveTime.Reset()
	m.prometheus.agentPoolJobFinishTime.Reset()

}

func (m *MetricsCollectorAgentPool) Collect(ctx context.Context, logger *log.Entry, callback chan<- func()) {
	for _, project := range m.CollectorReference.azureDevOpsProjects.List {
		contextLogger := logger.WithFields(log.Fields{
			"project": project.Name,
		})
		m.collectAgentInfo(ctx, contextLogger, callback, project)
	}

	for _, agentPoolId := range m.CollectorReference.AgentPoolIdList {
		contextLogger := logger.WithFields(log.Fields{
			"agentPoolId": agentPoolId,
		})

		m.collectAgentQueues(ctx, contextLogger, callback, agentPoolId)
		m.collectAgentPoolJobs(ctx, contextLogger, callback, agentPoolId)
	}
}

func (m *MetricsCollectorAgentPool) collectAgentInfo(ctx context.Context, logger *log.Entry, callback chan<- func(), project devopsClient.Project) {
	list, err := AzureDevopsClient.ListAgentQueues(project.Id)
	if err != nil {
		logger.Error(err)
		return
	}

	agentPoolInfoMetric := prometheusCommon.NewMetricsList()
	agentPoolSizeMetric := prometheusCommon.NewMetricsList()

	for _, agentQueue := range list.List {
		agentPoolInfoMetric.Add(prometheus.Labels{
			"agentPoolID":   int64ToString(agentQueue.Pool.Id),
			"agentPoolName": agentQueue.Name,
			"isHosted":      boolToString(agentQueue.Pool.IsHosted),
			"agentPoolType": agentQueue.Pool.PoolType,
		}, 1)

		agentPoolSizeMetric.Add(prometheus.Labels{
			"agentPoolID": int64ToString(agentQueue.Pool.Id),
		}, float64(agentQueue.Pool.Size))
	}

	callback <- func() {
		agentPoolInfoMetric.GaugeSet(m.prometheus.agentPool)
		agentPoolSizeMetric.GaugeSet(m.prometheus.agentPoolSize)
	}
}

func (m *MetricsCollectorAgentPool) collectAgentQueues(ctx context.Context, logger *log.Entry, callback chan<- func(), agentPoolId int64) {
	list, err := AzureDevopsClient.ListAgentPoolAgents(agentPoolId)
	if err != nil {
		logger.Error(err)
		return
	}

	agentPoolUsageMetric := prometheusCommon.NewMetricsList()
	agentPoolAgentMetric := prometheusCommon.NewMetricsList()
	agentPoolAgentStatusMetric := prometheusCommon.NewMetricsList()
	agentPoolAgentJobMetric := prometheusCommon.NewMetricsList()

	agentPoolSize := 0
	agentPoolUsed := 0

	logger.Infof("Traversing the agentpool agent list %v", list.Count)
	for _, agentPoolAgent := range list.List {
		agentPoolSize++
		infoLabels := prometheus.Labels{
			"agentPoolID":           int64ToString(agentPoolId),
			"agentPoolAgentID":      int64ToString(agentPoolAgent.Id),
			"agentPoolAgentName":    agentPoolAgent.Name,
			"agentPoolAgentVersion": agentPoolAgent.Version,
			"provisioningState":     agentPoolAgent.ProvisioningState,
			"maxParallelism":        int64ToString(agentPoolAgent.MaxParallelism),
			"agentPoolAgentOs":      agentPoolAgent.OsDescription,
			"enabled":               boolToString(agentPoolAgent.Enabled),
			"status":                agentPoolAgent.Status,
			"hasAssignedRequest":    boolToString(agentPoolAgent.AssignedRequest.RequestId > 0),
		}

		agentPoolAgentMetric.Add(infoLabels, 1)
		// logger.WithField("agentPoolAgentMetric:", agentPoolAgentMetric).Infof("These is the agentpoolAgent metric")

		statusCreatedLabels := prometheus.Labels{
			"agentPoolAgentID": int64ToString(agentPoolAgent.Id),
			"type":             "created",
		}
		agentPoolAgentStatusMetric.Add(statusCreatedLabels, timeToFloat64(agentPoolAgent.CreatedOn))

		if agentPoolAgent.AssignedRequest.RequestId > 0 {
			agentPoolUsed++
			jobLabels := prometheus.Labels{
				"agentPoolAgentID": int64ToString(agentPoolAgent.Id),
				"planType":         agentPoolAgent.AssignedRequest.PlanType,
				"jobRequestId":     int64ToString(agentPoolAgent.AssignedRequest.RequestId),
				"definitionID":     int64ToString(agentPoolAgent.AssignedRequest.Definition.Id),
				"definitionName":   agentPoolAgent.AssignedRequest.Definition.Name,
				"scopeID":          agentPoolAgent.AssignedRequest.ScopeId,
			}
			agentPoolAgentJobMetric.Add(jobLabels, timeToFloat64(agentPoolAgent.AssignedRequest.AssignTime))
			// logger.WithField("agentPoolJobInfoMetric:", agentPoolAgentJobMetric).Infof("These is the agentpoolAgentJobInfo metric")
		}
	}

	agentPoolUsageMetric.Add(prometheus.Labels{
		"agentPoolID": int64ToString(agentPoolId),
	}, float64(agentPoolUsed)/float64(agentPoolSize))

	callback <- func() {
		// TODO: loggerwithfield, log every metric
		agentPoolUsageMetric.GaugeSet(m.prometheus.agentPoolUsage)
		agentPoolAgentMetric.GaugeSet(m.prometheus.agentPoolAgent)
		agentPoolAgentStatusMetric.GaugeSet(m.prometheus.agentPoolAgentStatus)
		agentPoolAgentJobMetric.GaugeSet(m.prometheus.agentPoolAgentJob)
	}
}

func (m *MetricsCollectorAgentPool) collectAgentPoolJobs(ctx context.Context, logger *log.Entry, callback chan<- func(), agentPoolId int64) {
	list, err := AzureDevopsClient.ListAgentPoolJobs(agentPoolId)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.WithField("Number of agentpool jobs:", list.Count).Infof("This is the number of agentpool jobs in the pool %v", agentPoolId)

	agentPoolQueueLengthMetric := prometheusCommon.NewMetricsList()
	agentPoolJobRequestInfoMetric := prometheusCommon.NewMetricsList()
	agentPoolJobAssignTimeMetric := prometheusCommon.NewMetricsList()
	agentPoolJobQueueTimeMetric := prometheusCommon.NewMetricsList()
	agentPoolJobReceiveTimeMetric := prometheusCommon.NewMetricsList()
	agentPoolJobFinishTimeMetric := prometheusCommon.NewMetricsList()

	notStartedJobCount := 0

	for _, agentPoolJob := range list.List {

		if agentPoolJob.AssignTime.IsZero() {
			notStartedJobCount++
		}

		if timeToFloat64(agentPoolJob.FinishTime) > 0 && timeToFloat64(agentPoolJob.AssignTime) > 0 && agentPoolJob.Result != "canceled" {
			jobLabels := prometheus.Labels{
				"Name":         agentPoolJob.Definition.Name,
				"jobRequestID": int64ToString(agentPoolJob.RequestId),
				"JobID":        agentPoolJob.JobId,
				"Result":       agentPoolJob.Result,
			}

			// logger.WithField("infoLabels:", infoLabels).Infof("These are the joblabels added to agentpoolJobRequest metric")
			agentPoolJobRequestInfoMetric.AddInfo(jobLabels)
			// logger.WithField("agentPoolJobRequestMetric:", agentPoolJobRequestInfoMetric).Infof("These is the agentpoolAgentJob metric")

			agentPoolJobAssignTimeMetric.Add(prometheus.Labels{
				"jobRequestID": int64ToString(agentPoolJob.RequestId),
			}, timeToFloat64(agentPoolJob.AssignTime))

			agentPoolJobQueueTimeMetric.Add(prometheus.Labels{
				"jobRequestID": int64ToString(agentPoolJob.RequestId),
			}, timeToFloat64(agentPoolJob.QueueTime))

			agentPoolJobReceiveTimeMetric.Add(prometheus.Labels{
				"jobRequestID": int64ToString(agentPoolJob.RequestId),
			}, timeToFloat64(agentPoolJob.ReceiveTime))

			agentPoolJobFinishTimeMetric.Add(prometheus.Labels{
				"jobRequestID": int64ToString(agentPoolJob.RequestId),
			}, timeToFloat64(agentPoolJob.FinishTime))
		}

	}

	infoLabels := prometheus.Labels{
		"agentPoolID": int64ToString(agentPoolId),
	}

	agentPoolQueueLengthMetric.Add(infoLabels, float64(notStartedJobCount))

	callback <- func() {
		agentPoolQueueLengthMetric.GaugeSet(m.prometheus.agentPoolQueueLength)
		agentPoolJobRequestInfoMetric.GaugeSet(m.prometheus.agentPoolJobRequestInfo)
		agentPoolJobAssignTimeMetric.GaugeSet(m.prometheus.agentPoolJobAssignTime)
		agentPoolJobQueueTimeMetric.GaugeSet(m.prometheus.agentPoolJobQueueTime)
		agentPoolJobReceiveTimeMetric.GaugeSet(m.prometheus.agentPoolJobReceiveTime)
		agentPoolJobFinishTimeMetric.GaugeSet(m.prometheus.agentPoolJobFinishTime)
		// logger.WithField("agentPoolJobRequestGuage:", *m.prometheus.agentPoolJobRequest).Infof("These is the agentpoolAgentJobRequest prometheus gauge")
	}
}
